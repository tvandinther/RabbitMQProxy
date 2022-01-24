using RabbitMQ.Client;
using RabbitMQProxy;
using RabbitMQProxy.Services;

var builder = WebApplication.CreateBuilder(args);

var services = builder.Services;

services.AddScoped<IAmqpService, RabbitMqService>();

services.AddControllers();
services.AddEndpointsApiExplorer();
services.AddSwaggerGen();

var app = builder.Build();

app.UseSwagger();
app.UseSwaggerUI();

app.UseHttpsRedirection();

app.UseAuthorization();

app.MapControllers();

app.Run();