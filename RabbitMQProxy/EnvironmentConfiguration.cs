namespace RabbitMQProxy;

public static class EnvironmentConfiguration
{
    public static string BrokerHostName => Environment.GetEnvironmentVariable("BROKER_HOSTNAME") ?? "localhost";
    public static int BrokerPort => Environment.GetEnvironmentVariable("BROKER_PORT") != null ? int.Parse(Environment.GetEnvironmentVariable("RABBITMQ_PORT")) : 5672;
    public static string BrokerUserName => Environment.GetEnvironmentVariable("BROKER_USERNAME") ?? "guest";
    public static string BrokerPassword => Environment.GetEnvironmentVariable("BROKER_PASSWORD") ?? "guest";
}