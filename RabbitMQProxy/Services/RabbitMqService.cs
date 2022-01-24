using RabbitMQ.Client;

namespace RabbitMQProxy.Services;

public class RabbitMqService : IAmqpService
{
    private readonly ILogger<RabbitMqService> _logger;
    private readonly IConnectionFactory _connectionFactory;
    public string HostName { get; }
    public int Port { get; }

    public RabbitMqService(ILogger<RabbitMqService> logger)
    {
        _logger = logger;
        
        HostName = EnvironmentConfiguration.BrokerHostName;
        Port = EnvironmentConfiguration.BrokerPort;
        
        _connectionFactory = new ConnectionFactory
        {
            HostName = HostName,
            Port = Port,
            Password = EnvironmentConfiguration.BrokerPassword,
            UserName = EnvironmentConfiguration.BrokerUserName,
            ClientProvidedName = "proxy"
        };
    }

    public void SendMessage(string queueName, byte[] message)
    {
        using var connection = _connectionFactory.CreateConnection();
        using var channel = connection.CreateModel();

        channel.QueueDeclare(queueName);
        
        var properties = channel.CreateBasicProperties();
        
        channel.BasicPublish("", queueName, properties,message);
        
        _logger.LogInformation("Published message to queue {queueName}. Size: {size} bytes.", queueName, message.Length);
    }
}