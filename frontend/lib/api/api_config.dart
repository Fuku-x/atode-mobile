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
        ? 'https://api.example.com'
        : 'http://localhost:8080';

    final baseUrlRaw = const String.fromEnvironment(
      'API_BASE_URL',
      defaultValue: '',
    );

    final baseUrl = baseUrlRaw.isEmpty ? defaultBaseUrl : baseUrlRaw;

    final uri = Uri.parse(baseUrl);
    if (environment == ApiEnvironment.prod && uri.scheme != 'https') {
      throw StateError('In prod, API_BASE_URL must use https.');
    }

    return ApiConfig(
      environment: environment,
      baseUri: uri,
    );
  }
}
