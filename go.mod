module github.com/estafette/estafette-ci-api

go 1.14

require (
	cloud.google.com/go/bigquery v1.5.0
	cloud.google.com/go/pubsub v1.3.0
	cloud.google.com/go/storage v1.6.0
	github.com/Masterminds/squirrel v1.2.0
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ericchiang/k8s v1.2.0
	github.com/estafette/estafette-ci-contracts v0.0.182
	github.com/estafette/estafette-ci-crypt v0.0.36
	github.com/estafette/estafette-ci-manifest v0.1.145
	github.com/estafette/estafette-foundation v0.0.53
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gin-contrib/gzip v0.0.2-0.20190827144029-5602d8b438ea
	github.com/gin-gonic/gin v1.5.0
	github.com/go-kit/kit v0.10.0
	github.com/lib/pq v1.3.0
	github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.5.0
	github.com/rs/zerolog v1.18.0
	github.com/sethgrid/pester v0.0.0-20190127155807-68a33a018ad0
	github.com/stretchr/testify v1.5.1
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	google.golang.org/api v0.20.0
	gopkg.in/yaml.v2 v2.2.8
)
