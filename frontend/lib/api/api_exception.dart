enum ApiExceptionType {
  http,
  network,
  timeout,
  invalidResponse,
}

class ApiException implements Exception {
  ApiException({
    required this.type,
    required this.message,
    this.uri,
    this.statusCode,
    this.body,
    this.cause,
  });

  final ApiExceptionType type;
  final String message;
  final Uri? uri;
  final int? statusCode;
  final dynamic body;
  final Object? cause;

  factory ApiException.http({
    required int statusCode,
    required String message,
    Uri? uri,
    dynamic body,
    Object? cause,
  }) {
    return ApiException(
      type: ApiExceptionType.http,
      statusCode: statusCode,
      message: message,
      uri: uri,
      body: body,
      cause: cause,
    );
  }

  factory ApiException.network({
    required String message,
    Uri? uri,
    Object? cause,
  }) {
    return ApiException(
      type: ApiExceptionType.network,
      message: message,
      uri: uri,
      cause: cause,
    );
  }

  factory ApiException.timeout({
    required String message,
    Uri? uri,
    Object? cause,
  }) {
    return ApiException(
      type: ApiExceptionType.timeout,
      message: message,
      uri: uri,
      cause: cause,
    );
  }

  factory ApiException.invalidResponse({
    required String message,
    Uri? uri,
    Object? cause,
  }) {
    return ApiException(
      type: ApiExceptionType.invalidResponse,
      message: message,
      uri: uri,
      cause: cause,
    );
  }

  @override
  String toString() {
    final code = statusCode == null ? '' : ' statusCode=$statusCode';
    final u = uri == null ? '' : ' uri=$uri';
    return 'ApiException(type=$type$message$code$u)';
  }
}
