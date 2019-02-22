mustOutputDir=$1
if test -z "$mustOutputDir"; then
	echo "Error! Output dir must be defined."
	return
fi

mustTag=$(git tag -l --points-at HEAD)
res=""

if test -z "$mustTag"; then
	res=$(git rev-parse HEAD)
else
	res=$mustTag
fi

go build -ldflags "-X main.version=$res" -o /usr/local/bin/horizon gitlab.com/tokend/horizon/cmd/horizon
