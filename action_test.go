package horizon

import (
	"database/sql"
	"testing"

	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/log"
	"gitlab.com/tokend/horizon/render/problem"
	"github.com/go-errors/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIsAllowedAction(t *testing.T) {
	Convey("Given accountShowAction", t, func() {
		dataOwner, _ := keypair.Random()
		coreQ := core.CoreQMock{}
		action := Action{
			App: NewTestApp(),
			cq:  &coreQ,
			hq:  &history.QMock{},
			Log: log.WithField("testing", "action:is_allowed"),
		}
		//action.App.CoreInfo.MasterAccountID = "master"
		accountShowAction := AccountShowAction{
			Action:  action,
			Address: dataOwner.Address(),
		}
		Convey("Not signed request", func() {
			accountShowAction.IsAllowed(dataOwner.Address())
			assert.Equal(t, problem.SignNotVerified.Error(), accountShowAction.Err.Error())
		})

		accountShowAction.IsSigned = true
		masterAccountID := action.App.CoreInfo.MasterAccountID
		Convey("random db error while getting master account", func() {
			requestSigner, _ := keypair.Random()
			accountShowAction.Signer = requestSigner.Address()
			account := &core.Account{}
			coreQ.On("AccountByAddress", account, masterAccountID).Return(errors.New("something wrong")).Once()
			accountShowAction.IsAllowed(dataOwner.Address())
			assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
		})
		Convey("master found in db", func() {
			account := &core.Account{}
			Convey("requestSigner is master", func() {
				accountShowAction.Signer = masterAccountID
				Convey("master has Thresholds[0] = 0", func() {
					coreQ.On("AccountByAddress", account, masterAccountID).Run(func(args mock.Arguments) {
						account := args.Get(0).(*core.Account)
						account.Accountid = masterAccountID
						account.AccountType = int32(xdr.AccountTypeMaster)
						account.Thresholds[0] = 0
					}).Return(nil).Once()
					Convey("failure to get signers of master", func() {
						var signers []core.Signer
						coreQ.On("SignersByAddress", &signers, masterAccountID).Return(sql.ErrNoRows).Once()
						accountShowAction.IsAllowed(dataOwner.Address())
						assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
					})
					Convey("got signers of master", func() {
						var signers []core.Signer
						// Content does not matter because thresholds[0] is 0
						coreQ.On("SignersByAddress", &signers, masterAccountID).Return(nil).Once()

						Convey("empty owner of data", func() {
							accountShowAction.IsAllowed("")
							assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
						})
						Convey("data owner is master", func() {
							accountShowAction.IsAllowed(masterAccountID)
							assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
						})
						Convey("error while obtaining dataOwner", func() {
							coreQ.On("AccountByAddress", account, dataOwner.Address()).Return(errors.New("something wrong")).Once()
							accountShowAction.IsAllowed(dataOwner.Address())
							assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
						})
						Convey("dataOwner not found in db", func() {
							coreQ.On("AccountByAddress", account, dataOwner.Address()).Return(sql.ErrNoRows).Once()
							accountShowAction.IsAllowed(dataOwner.Address())
							assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
						})
						Convey("dataOwner account found in db and is general user", func() {
							coreQ.On("AccountByAddress", account, dataOwner.Address()).Run(func(args mock.Arguments) {
								account := args.Get(0).(*core.Account)
								account.Accountid = dataOwner.Address()
								account.AccountType = int32(xdr.AccountTypeGeneral)
							}).Return(nil).Once()
							Convey("error while obtaining signers of dataOwner", func() {
								var signers []core.Signer
								coreQ.On("SignersByAddress", &signers, dataOwner.Address()).Return(sql.ErrNoRows).Once()
								accountShowAction.IsAllowed(dataOwner.Address())
								assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
							})
							Convey("got requestSigner among signers of dataOwner", func() {
								var signers []core.Signer
								coreQ.On("SignersByAddress", &signers, dataOwner.Address()).Run(func(args mock.Arguments) {
									signers := args.Get(0).(*[]core.Signer)
									signer := core.Signer{
										Publickey: accountShowAction.Signer,
									}
									*signers = append(*signers, signer)
								}).Return(nil).Once()
								accountShowAction.IsAllowed(dataOwner.Address())
								assert.Equal(t, nil, accountShowAction.Err)
							})
							Convey("no requestSigner among signers of dataOwner", func() {
								var signers []core.Signer
								coreQ.On("SignersByAddress", &signers, dataOwner.Address()).Return(nil).Once()
								accountShowAction.IsAllowed(dataOwner.Address())
								assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
							})
						})
						Convey("dataOwner is comissionAccount or sequenceProvider found in db", func() {
							coreQ.On("AccountByAddress", account, dataOwner.Address()).Run(func(args mock.Arguments) {
								account := args.Get(0).(*core.Account)
								account.Accountid = dataOwner.Address()
								account.AccountType = int32(xdr.AccountTypeCommission)
							}).Return(nil).Once()
							Convey("master not loaded", func() {
								masterAcc := &core.Account{}
								coreQ.On("AccountByAddress", masterAcc, masterAccountID).Return(sql.ErrNoRows).Once()
								accountShowAction.IsAllowed(dataOwner.Address())
								assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
							})
							Convey("master loaded", func() {
								masterAcc := &core.Account{}
								coreQ.On("AccountByAddress", masterAcc, masterAccountID).Run(func(args mock.Arguments) {
									account := args.Get(0).(*core.Account)
									account.Accountid = masterAccountID
									account.AccountType = int32(xdr.AccountTypeMaster)
								}).Return(nil).Once()
								Convey("master signers not loaded", func() {
									var signers []core.Signer
									coreQ.On("SignersByAddress", &signers, masterAccountID).Return(sql.ErrNoRows).Once()
									accountShowAction.IsAllowed(dataOwner.Address())
									assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
								})
								Convey("master signers loaded and contains requestSigner", func() {
									// Impossible because requestSigner == master
								})
								Convey("master signers loaded and does not contain requestSigner", func() {
									coreQ.On("SignersByAddress", &signers, masterAccountID).Run(func(args mock.Arguments) {
										randomAccount, _ := keypair.Random()
										signers := args.Get(0).(*[]core.Signer)
										signer := core.Signer{
											Publickey: randomAccount.Address(),
										}
										*signers = append(*signers, signer)
									}).Return(nil).Once()
									accountShowAction.IsAllowed(dataOwner.Address())
									assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
								})
							})

						})
					})
				})
				Convey("master has Thresholds[0] = 1", func() {
					coreQ.On("AccountByAddress", account, masterAccountID).Run(func(args mock.Arguments) {
						account := args.Get(0).(*core.Account)
						account.Accountid = masterAccountID
						account.AccountType = int32(xdr.AccountTypeMaster)
						account.Thresholds[0] = 1
					}).Return(nil).Once()
					Convey("failure to get signers of master", func() {
						var signers []core.Signer
						coreQ.On("SignersByAddress", &signers, masterAccountID).Return(sql.ErrNoRows).Once()
						accountShowAction.IsAllowed(dataOwner.Address())
						assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
					})
					Convey("got signers of master", func() {
						var signers []core.Signer
						coreQ.On("SignersByAddress", &signers, masterAccountID).Return(nil).Once()
						accountShowAction.IsAllowed(dataOwner.Address())
						assert.Equal(t, nil, accountShowAction.Err)
					})
				})
			})

			coreQ.On("AccountByAddress", account, masterAccountID).Run(func(args mock.Arguments) {
				account := args.Get(0).(*core.Account)
				account.Accountid = masterAccountID
				account.AccountType = int32(xdr.AccountTypeMaster)
			}).Return(nil).Once()

			Convey("requestSigner is not master", func() {
				requestSigner, _ := keypair.Random()
				accountShowAction.Signer = requestSigner.Address()

				Convey("failure to get signers of master", func() {
					var signers []core.Signer
					coreQ.On("SignersByAddress", &signers, masterAccountID).Return(sql.ErrNoRows).Once()
					accountShowAction.IsAllowed(dataOwner.Address())
					assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
				})
				Convey("requestSigner among signers of master", func() {
					var signers []core.Signer
					coreQ.On("SignersByAddress", &signers, masterAccountID).Run(func(args mock.Arguments) {
						signers := args.Get(0).(*[]core.Signer)
						signer := core.Signer{
							Publickey: requestSigner.Address(),
						}
						*signers = append(*signers, signer)
					}).Return(nil).Once()
					accountShowAction.IsAllowed("")
					assert.Equal(t, nil, accountShowAction.Err)
				})
				Convey("no requestSigner among signers of master", func() {
					var signers []core.Signer
					coreQ.On("SignersByAddress", &signers, masterAccountID).Run(func(args mock.Arguments) {
						randomAccount, _ := keypair.Random()
						signers := args.Get(0).(*[]core.Signer)
						signer := core.Signer{
							Publickey: randomAccount.Address(),
						}
						*signers = append(*signers, signer)
					}).Return(nil).Once()

					Convey("empty owner", func() {
						accountShowAction.IsAllowed("")
						assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
					})
					Convey("owner is master", func() {
						accountShowAction.IsAllowed(masterAccountID)
						assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
					})
					Convey("error while obtaining dataOwner", func() {
						coreQ.On("AccountByAddress", account, dataOwner.Address()).Return(errors.New("something wrong")).Once()
						accountShowAction.IsAllowed(dataOwner.Address())
						assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
					})
					Convey("dataOwner not found in db", func() {
						coreQ.On("AccountByAddress", account, dataOwner.Address()).Return(sql.ErrNoRows).Once()
						accountShowAction.IsAllowed(dataOwner.Address())
						assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
					})
					Convey("general dataOwner account found in db", func() {
						coreQ.On("AccountByAddress", account, dataOwner.Address()).Run(func(args mock.Arguments) {
							account := args.Get(0).(*core.Account)
							account.Accountid = dataOwner.Address()
							account.AccountType = int32(xdr.AccountTypeGeneral)
						}).Return(nil).Once()
						Convey("failure to get signers of dataOwner", func() {
							var signers []core.Signer
							coreQ.On("SignersByAddress", &signers, dataOwner.Address()).Return(sql.ErrNoRows).Once()
							accountShowAction.IsAllowed(dataOwner.Address())
							assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
						})
						Convey("no requestSigner among signers of dataOwner", func() {
							var signers []core.Signer
							coreQ.On("SignersByAddress", &signers, dataOwner.Address()).Run(func(args mock.Arguments) {
								randomAccount, _ := keypair.Random()
								signers := args.Get(0).(*[]core.Signer)
								signer := core.Signer{
									Publickey: randomAccount.Address(),
								}
								*signers = append(*signers, signer)
							}).Return(nil).Once()
							accountShowAction.IsAllowed(dataOwner.Address())
							assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
						})
						Convey("got requestSigner among signers of dataOwners", func() {
							var signers []core.Signer
							coreQ.On("SignersByAddress", &signers, dataOwner.Address()).Run(func(args mock.Arguments) {
								signers := args.Get(0).(*[]core.Signer)
								signer := core.Signer{
									Publickey: requestSigner.Address(),
								}
								*signers = append(*signers, signer)
							}).Return(nil).Once()
							accountShowAction.IsAllowed(dataOwner.Address())
							assert.Equal(t, nil, accountShowAction.Err)
						})
					})
					Convey("dataOwner is comissionAccount or sequenceProvider found in db", func() {
						coreQ.On("AccountByAddress", account, dataOwner.Address()).Run(func(args mock.Arguments) {
							account := args.Get(0).(*core.Account)
							account.Accountid = dataOwner.Address()
							account.AccountType = int32(xdr.AccountTypeCommission)
						}).Return(nil).Once()
						Convey("master not loaded", func() {
							masterAcc := &core.Account{}
							coreQ.On("AccountByAddress", masterAcc, masterAccountID).Return(sql.ErrNoRows).Once()
							accountShowAction.IsAllowed(dataOwner.Address())
							assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
						})
						Convey("master loaded", func() {
							masterAcc := &core.Account{}
							coreQ.On("AccountByAddress", masterAcc, masterAccountID).Run(func(args mock.Arguments) {
								account := args.Get(0).(*core.Account)
								account.Accountid = masterAccountID
								account.AccountType = int32(xdr.AccountTypeMaster)
							}).Return(nil).Once()
							var signers []core.Signer
							Convey("master signers not loaded", func() {
								coreQ.On("SignersByAddress", &signers, masterAccountID).Return(sql.ErrNoRows).Once()
								accountShowAction.IsAllowed(dataOwner.Address())
								assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
							})
							Convey("master signers loaded and contains requestSigners", func() {
								coreQ.On("SignersByAddress", &signers, masterAccountID).Run(func(args mock.Arguments) {
									signers := args.Get(0).(*[]core.Signer)
									signer := core.Signer{
										Publickey: requestSigner.Address(),
									}
									*signers = append(*signers, signer)
								}).Return(nil).Once()
								accountShowAction.IsAllowed(dataOwner.Address())
								assert.Equal(t, nil, accountShowAction.Err)
							})
							Convey("master signers loaded and does not contain requestSigners", func() {
								coreQ.On("SignersByAddress", &signers, masterAccountID).Run(func(args mock.Arguments) {
									randomAccount, _ := keypair.Random()
									signers := args.Get(0).(*[]core.Signer)
									signer := core.Signer{
										Publickey: randomAccount.Address(),
									}
									*signers = append(*signers, signer)
								}).Return(nil).Once()
								accountShowAction.IsAllowed(dataOwner.Address())
								assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
							})
						})
					})

				})
			})

		})
		Convey("master not found in db", func() {
			requestSigner, _ := keypair.Random()
			accountShowAction.Signer = requestSigner.Address()
			account := &core.Account{}
			coreQ.On("AccountByAddress", account, masterAccountID).Return(sql.ErrNoRows).Once()
			Convey("empty owner", func() {
				accountShowAction.IsAllowed("")
				assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
			})
			Convey("owner is master", func() {
				accountShowAction.IsAllowed(masterAccountID)
				assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
			})
			Convey("error while obtaining dataOwner", func() {
				coreQ.On("AccountByAddress", account, dataOwner.Address()).Return(errors.New("something wrong")).Once()
				accountShowAction.IsAllowed(dataOwner.Address())
				assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
			})
			Convey("dataOwner not found in db", func() {
				coreQ.On("AccountByAddress", account, dataOwner.Address()).Return(sql.ErrNoRows).Once()
				accountShowAction.IsAllowed(dataOwner.Address())
				assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
			})
			Convey("general dataOwner account found in db", func() {
				coreQ.On("AccountByAddress", account, dataOwner.Address()).Run(func(args mock.Arguments) {
					account := args.Get(0).(*core.Account)
					account.Accountid = dataOwner.Address()
					account.AccountType = int32(xdr.AccountTypeGeneral)
				}).Return(nil).Once()
				Convey("failure to get signers of dataOwner", func() {
					var signers []core.Signer
					coreQ.On("SignersByAddress", &signers, dataOwner.Address()).Return(sql.ErrNoRows).Once()
					accountShowAction.IsAllowed(dataOwner.Address())
					assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
				})
				Convey("no requestSigner among signers of dataOwner", func() {
					var signers []core.Signer
					coreQ.On("SignersByAddress", &signers, dataOwner.Address()).Run(func(args mock.Arguments) {
						randomAccount, _ := keypair.Random()
						signers := args.Get(0).(*[]core.Signer)
						signer := core.Signer{
							Publickey: randomAccount.Address(),
						}
						*signers = append(*signers, signer)
					}).Return(nil).Once()
					accountShowAction.IsAllowed(dataOwner.Address())
					assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
				})
				Convey("got requestSigner among signers of dataOwners", func() {
					var signers []core.Signer
					coreQ.On("SignersByAddress", &signers, dataOwner.Address()).Run(func(args mock.Arguments) {
						signers := args.Get(0).(*[]core.Signer)
						signer := core.Signer{
							Publickey: requestSigner.Address(),
						}
						*signers = append(*signers, signer)
					}).Return(nil).Once()
					accountShowAction.IsAllowed(dataOwner.Address())
					assert.Equal(t, nil, accountShowAction.Err)
				})
			})
			Convey("dataOwner is comissionAccount or sequenceProvider found in db", func() {
				coreQ.On("AccountByAddress", account, dataOwner.Address()).Run(func(args mock.Arguments) {
					account := args.Get(0).(*core.Account)
					account.Accountid = dataOwner.Address()
					account.AccountType = int32(xdr.AccountTypeCommission)
				}).Return(nil).Once()
				Convey("master not loaded", func() {
					masterAcc := &core.Account{}
					coreQ.On("AccountByAddress", masterAcc, masterAccountID).Return(sql.ErrNoRows).Once()
					accountShowAction.IsAllowed(dataOwner.Address())
					assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
				})
				Convey("master loaded", func() {
					masterAcc := &core.Account{}
					coreQ.On("AccountByAddress", masterAcc, masterAccountID).Run(func(args mock.Arguments) {
						account := args.Get(0).(*core.Account)
						account.Accountid = masterAccountID
						account.AccountType = int32(xdr.AccountTypeMaster)
					}).Return(nil).Once()
					var signers []core.Signer
					Convey("master signers not loaded", func() {
						coreQ.On("SignersByAddress", &signers, masterAccountID).Return(sql.ErrNoRows).Once()
						accountShowAction.IsAllowed(dataOwner.Address())
						assert.Equal(t, problem.ServerError.Error(), accountShowAction.Err.Error())
					})
					Convey("master signers loaded and contains requestSigners", func() {
						coreQ.On("SignersByAddress", &signers, masterAccountID).Run(func(args mock.Arguments) {
							signers := args.Get(0).(*[]core.Signer)
							signer := core.Signer{
								Publickey: requestSigner.Address(),
							}
							*signers = append(*signers, signer)
						}).Return(nil).Once()
						accountShowAction.IsAllowed(dataOwner.Address())
						assert.Equal(t, nil, accountShowAction.Err)
					})
					Convey("master signers loaded and does not contain requestSigners", func() {
						coreQ.On("SignersByAddress", &signers, masterAccountID).Run(func(args mock.Arguments) {
							randomAccount, _ := keypair.Random()
							signers := args.Get(0).(*[]core.Signer)
							signer := core.Signer{
								Publickey: randomAccount.Address(),
							}
							*signers = append(*signers, signer)
						}).Return(nil).Once()
						accountShowAction.IsAllowed(dataOwner.Address())
						assert.Equal(t, problem.NotAllowed.Error(), accountShowAction.Err.Error())
					})
				})

			})
		})
	})
}
