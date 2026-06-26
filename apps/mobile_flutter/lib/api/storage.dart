import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class Storage {
  Storage._();

  static const _storage = FlutterSecureStorage();

  static const _tokenKey = 'cafeos_token';
  static const _tenantIdKey = 'cafeos_tenant_id';

  static Future<String?> getToken() => _storage.read(key: _tokenKey);
  static Future<void> setToken(String token) =>
      _storage.write(key: _tokenKey, value: token);
  static Future<void> clearToken() => _storage.delete(key: _tokenKey);

  static Future<String?> getTenantId() => _storage.read(key: _tenantIdKey);
  static Future<void> setTenantId(String id) =>
      _storage.write(key: _tenantIdKey, value: id);
  static Future<void> clearTenantId() => _storage.delete(key: _tenantIdKey);

  static Future<void> clearAll() async {
    await _storage.deleteAll();
  }
}
