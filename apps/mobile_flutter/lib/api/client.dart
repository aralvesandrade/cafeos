import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'storage.dart';

class ApiClient {
  ApiClient._();

  static final Dio _dio = Dio(BaseOptions(
    baseUrl: const String.fromEnvironment(
      'API_BASE_URL',
      defaultValue: 'http://localhost:5001',
    ),
    connectTimeout: const Duration(seconds: 10),
    receiveTimeout: const Duration(seconds: 10),
    headers: {'Content-Type': 'application/json'},
  ));

  static Future<Response> request(
    String path, {
    String method = 'GET',
    dynamic body,
    Map<String, dynamic>? queryParams,
  }) async {
    final token = await Storage.getToken();
    final tenantId = await Storage.getTenantId();

    String url = path;
    if (token != null && tenantId != null && !path.startsWith('/auth/')) {
      url = '/api/v1/$tenantId$path';
    }

    final options = Options(
      method: method,
      headers: {
        if (token != null) 'Authorization': 'Bearer $token',
      },
    );

    try {
      final response = await _dio.request(
        url,
        data: body,
        queryParameters: queryParams,
        options: options,
      );
      return response;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  static Exception _handleError(DioException e) {
    if (e.response?.data is Map) {
      final msg = (e.response!.data as Map)['error'] as String?;
      if (msg != null) return Exception(msg);
    }
    return Exception('Erro de conexão com o servidor');
  }
}

final apiClientProvider = Provider<ApiClient>((ref) => ApiClient._());
