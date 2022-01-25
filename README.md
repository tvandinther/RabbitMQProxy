# rabbitmq-proxy
A proxy API for interfacing with RabbitMQ. Useful for integrating test pipelines with existing API tools.

The container repository can be found under [tvandinther/rabbitmq-proxy](https://hub.docker.com/r/tvandinther/rabbitmq-proxy)

## Usage

Send your message payload as the body in an HTTP Post request to `/queue/{queueName}`. The route supports the following query parameters to control the queue declaration.
`durable`, `autoDelete`, `exclusive` and `noWait`. To set these parameters add them to the url: `/queue/{queueName}?durable=true&exclusive=true`.

Parameter default values:

| Parameter  	| default 	|
|------------	|---------	|
| durable    	| false   	|
| autoDelete 	| false   	|
| exclusive  	| false   	|
| noWait     	| false   	|

RabbitMQ connectivity can be configured using the following environment variables:

| Variable          	| default   	|
|-------------------	|-----------	|
| RABBITMQ_HOST     	| localhost 	|
| RABBITMQ_PORT     	| 5672      	|
| RABBITMQ_USERNAME 	| guest     	|
| RABBITMQ_PASSWORD 	| guest     	|

## Notes

This proxy server currently only supports publishing to a queue on the default exchange. It also does not currently support headers and properties. I welcome pull requests to extend its functionality.
