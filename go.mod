module gitlab.com/tokend/horizon

go 1.10

replace gopkg.in/throttled/throttled.v2 => github.com/throttled/throttled/v2 v2.6.0

require (
	github.com/Masterminds/squirrel v1.4.0
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/certifi/gocertifi v0.0.0-20200922220541-2c3bb06c6054 // indirect
	github.com/cheekybits/genny v1.0.0
	github.com/evalphobia/logrus_sentry v0.8.2
	github.com/getsentry/raven-go v0.2.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-errors/errors v1.1.1
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/goji/context v0.0.0-20160122015720-68b83f7b0439
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e
	github.com/google/jsonapi v0.0.0-20200825183604-3e3da1210d0c
	github.com/guregu/null v3.2.0+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.8.0
	github.com/magiconair/properties v1.8.4
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
	github.com/nullstyle/go-xdr v0.0.0-20180726165426-f4c839f75077 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
	github.com/rs/cors v1.7.0
	github.com/rubenv/sql-migrate v0.0.0-20200616145509-8d140a17f351
	github.com/sirupsen/logrus v1.7.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/zenazn/goji v1.0.1
	gitlab.com/distributed_lab/ape v1.5.0
	gitlab.com/distributed_lab/corer v0.0.0-20171130114150-cbfb46525895
	gitlab.com/distributed_lab/figure v2.1.0+incompatible
	gitlab.com/distributed_lab/kit v1.8.1
	gitlab.com/distributed_lab/logan v3.7.2+incompatible
	gitlab.com/distributed_lab/lorem v0.2.0 // indirect
	gitlab.com/distributed_lab/running v0.0.0-20200706131153-4af0e83eb96c
	gitlab.com/distributed_lab/txsub v0.0.0-20171130120140-d7781cbc2319
	gitlab.com/distributed_lab/urlval v2.2.0+incompatible
	gitlab.com/tokend/go v3.13.1-0.20200928125835-32df2e4a6022+incompatible
	gitlab.com/tokend/keypair v0.0.0-20190412110653-b9d7e0c8b312 // indirect
	gitlab.com/tokend/regources v4.9.2-0.20200918161150-47835a98daca+incompatible
	golang.org/x/net v0.0.0-20200927032502-5d4f70055728
	gopkg.in/throttled/throttled.v2 v2.6.0
	gopkg.in/tylerb/graceful.v1 v1.2.15
)
