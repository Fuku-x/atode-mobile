import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'api_config.dart';
import 'api_exception.dart';

typedef TokenProvider = Future<String?> Function();

typedef ApiDecoder<T> = T Function(dynamic json);

class ApiClient {
  ApiClient({
    required ApiConfig config,
    HttpClient? httpClient,
    this.tokenProvider,
    Duration timeout = const Duration(seconds: 30),
    Map<String, String> defaultHeaders = const {},
  })  : _config = config,
        _httpClient = httpClient ?? HttpClient(),
        _timeout = timeout,
        _defaultHeaders = defaultHeaders;

  final ApiConfig _config;
  final HttpClient _httpClient;
  final Duration _timeout;
  final Map<String, String> _defaultHeaders;
  final TokenProvider? tokenProvider;

  Future<T> get<T>(
    String path, {
    Map<String, String>? queryParameters,
    Map<String, String>? headers,
    ApiDecoder<T>? decoder,
  }) {
    return _request<T>(
      'GET',
      path,
      queryParameters: queryParameters,
      headers: headers,
      decoder: decoder,
    );
  }

  Future<T> post<T>(
    String path, {
    Object? body,
    Map<String, String>? queryParameters,
    Map<String, String>? headers,
    ApiDecoder<T>? decoder,
  }) {
    return _request<T>(
      'POST',
      path,
      body: body,
      queryParameters: queryParameters,
      headers: headers,
      decoder: decoder,
    );
  }

  Future<T> put<T>(
    String path, {
    Object? body,
    Map<String, String>? queryParameters,
    Map<String, String>? headers,
    ApiDecoder<T>? decoder,
  }) {
    return _request<T>(
      'PUT',
      path,
      body: body,
      queryParameters: queryParameters,
      headers: headers,
      decoder: decoder,
    );
  }

  Future<T> delete<T>(
    String path, {
    Object? body,
    Map<String, String>? queryParameters,
    Map<String, String>? headers,
    ApiDecoder<T>? decoder,
  }) {
    return _request<T>(
      'DELETE',
      path,
      body: body,
      queryParameters: queryParameters,
      headers: headers,
      decoder: decoder,
    );
  }

  Future<T> _request<T>(
    String method,
    String path, {
    Object? body,
    Map<String, String>? queryParameters,
    Map<String, String>? headers,
    ApiDecoder<T>? decoder,
  }) async {
    final uri = _resolveUri(path, queryParameters);

    try {
      final request = await _openRequest(method, uri);

      final mergedHeaders = <String, String>{
        ..._defaultHeaders,
        ...?headers,
      };

      final token = tokenProvider == null ? null : await tokenProvider!();
      if (token != null && token.isNotEmpty) {
        mergedHeaders.putIfAbsent('Authorization', () => 'Bearer $token');
      }

      mergedHeaders.forEach((key, value) {
        request.headers.set(key, value);
      });

      if (body != null) {
        request.headers.contentType = ContentType.json;
        request.write(jsonEncode(body));
      }

      final response = await request.close().timeout(_timeout);
      final responseBodyText = await utf8.decoder.bind(response).join();

      final decodedBody = _decodeBody(response, responseBodyText);

      if (response.statusCode < 200 || response.statusCode >= 300) {
        throw ApiException.http(
          statusCode: response.statusCode,
          uri: uri,
          body: decodedBody,
          message: _errorMessageFromBody(decodedBody) ?? 'Request failed',
        );
      }

      final value = decodedBody;

      if (decoder != null) {
        return decoder(value);
      }

      return value as T;
    } on TimeoutException catch (e) {
      throw ApiException.timeout(
        message: 'Request timed out',
        uri: uri,
        cause: e,
      );
    } on SocketException catch (e) {
      throw ApiException.network(
        message: 'Network error',
        uri: uri,
        cause: e,
      );
    } on HandshakeException catch (e) {
      throw ApiException.network(
        message: 'TLS handshake failed',
        uri: uri,
        cause: e,
      );
    } on ApiException {
      rethrow;
    } catch (e) {
      throw ApiException.invalidResponse(
        message: 'Unexpected error',
        uri: uri,
        cause: e,
      );
    }
  }

  Future<HttpClientRequest> _openRequest(String method, Uri uri) {
    switch (method) {
      case 'GET':
        return _httpClient.getUrl(uri);
      case 'POST':
        return _httpClient.postUrl(uri);
      case 'PUT':
        return _httpClient.putUrl(uri);
      case 'DELETE':
        return _httpClient.deleteUrl(uri);
      default:
        return _httpClient.openUrl(method, uri);
    }
  }

  Uri _resolveUri(String path, Map<String, String>? queryParameters) {
    final base = _config.baseUri;

    final basePath = base.path;
    final nextPath = _joinPaths(basePath, path);

    return base.replace(
      path: nextPath,
      queryParameters: {
        ...base.queryParameters,
        ...?queryParameters,
      },
    );
  }

  String _joinPaths(String basePath, String path) {
    final b = basePath.isEmpty ? '/' : basePath;
    final normalizedBase = b.endsWith('/') ? b.substring(0, b.length - 1) : b;
    final normalizedPath = path.startsWith('/') ? path : '/$path';
    final joined = '$normalizedBase$normalizedPath';
    return joined.isEmpty ? '/' : joined;
  }

  dynamic _decodeBody(HttpClientResponse response, String bodyText) {
    if (bodyText.isEmpty) return null;

    final contentType = response.headers.contentType?.mimeType ??
        response.headers.value(HttpHeaders.contentTypeHeader) ??
        '';

    final isJson = contentType.contains('application/json') ||
        contentType.contains('+json') ||
        contentType.contains('text/json');

    if (!isJson) return bodyText;

    try {
      return jsonDecode(bodyText);
    } catch (e) {
      throw ApiException.invalidResponse(
        message: 'Failed to decode JSON response',
        cause: e,
      );
    }
  }

  String? _errorMessageFromBody(dynamic body) {
    if (body is Map) {
      final message = body['message'];
      if (message is String && message.isNotEmpty) return message;
    }
    return null;
  }
}
