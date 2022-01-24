FROM mcr.microsoft.com/dotnet/aspnet:6.0 AS base
WORKDIR /app
EXPOSE 80

FROM mcr.microsoft.com/dotnet/sdk:6.0 AS build
WORKDIR /src
COPY ["RabbitMQProxy/RabbitMQProxy.csproj", "RabbitMQProxy/"]
RUN dotnet restore "RabbitMQProxy/RabbitMQProxy.csproj"
COPY . .
WORKDIR "/src/RabbitMQProxy"
RUN dotnet build "RabbitMQProxy.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "RabbitMQProxy.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "RabbitMQProxy.dll"]
