enum ApiEnvironment { dev, prod }

class ApiConfig {
  ApiConfig({required this.environment, required this.baseUri});

  final ApiEnvironment environment;
  final Uri baseUri;

  factory ApiConfig.fromEnvironment() {
    final envRaw = const String.fromEnvironment('API_ENV', defaultValue: 'dev');
    final environment = envRaw.toLowerCase() == 'prod'
        ? ApiEnvironment.prod
        : ApiEnvironment.dev;

    final defaultBaseUrl = environment == ApiEnvironment.prod
        ? 'http://localhost:8080'
        : 'http://localhost:8080';

    final baseUrlRaw = const String.fromEnvironment(
      'API_BASE_URL',
      defaultValue: '',
    );

    final baseUrl = baseUrlRaw.isEmpty ? defaultBaseUrl : baseUrlRaw;

    return ApiConfig(
      environment: environment,
      baseUri: Uri.parse(baseUrl),
    );
  }
}
