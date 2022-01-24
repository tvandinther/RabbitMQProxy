namespace RabbitMQProxy.Services;

public interface IAmqpService
{
    public string HostName { get; }
    public int Port { get; }
    
    void SendMessage(string queueName, byte[] message);
}