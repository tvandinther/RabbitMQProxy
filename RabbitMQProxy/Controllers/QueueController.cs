using System.Net;
using Microsoft.AspNetCore.Mvc;
using RabbitMQ.Client.Exceptions;
using RabbitMQProxy.Services;

namespace RabbitMQProxy.Controllers;

[ApiController]
[Route("[controller]")]
public class QueueController : ControllerBase
{
    private readonly ILogger<QueueController> _logger;
    private readonly IAmqpService _amqpService;

    public QueueController(ILogger<QueueController> logger, IAmqpService amqpService)
    {
        _logger = logger;
        _amqpService = amqpService;
    }
    
    [HttpPost("{queueName}")]
    public async Task<IActionResult> Post([FromRoute] string queueName)
    {
        await using var ms = new MemoryStream();
        await Request.Body.CopyToAsync(ms);
        var message = ms.ToArray();

        try
        {
            _amqpService.SendMessage(queueName, message);
        }
        catch (BrokerUnreachableException e)
        {
            var errorMessage = String.Format("Broker unreachable at {0}:{1}", _amqpService.HostName,
                _amqpService.Port);
            _logger.LogWarning(e, errorMessage);

            return StatusCode((int) HttpStatusCode.BadGateway, errorMessage);
        }

        return Ok();
    }
}