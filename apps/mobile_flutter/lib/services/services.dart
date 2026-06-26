import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:uuid/uuid.dart';
import '../api/client.dart';
import '../api/storage.dart';
import '../db/schema.dart' as db;
import '../models/models.dart' as models;
import '../repos/operation_repo.dart';
import '../repos/plot_repo.dart';
import '../repos/sync_queue_repo.dart';
import 'sync_service.dart';

class AuthService {
  Future<AuthResult> login(String email, String password) async {
    final response = await ApiClient.request(
      '/auth/login',
      method: 'POST',
      body: {'email': email, 'password': password},
    );

    final data = response.data as Map<String, dynamic>;
    final token = data['token'] as String;
    final tenantId = data['tenant_id'] as String;

    await Storage.setToken(token);
    await Storage.setTenantId(tenantId);

    return AuthResult(
      token: token,
      tenantId: tenantId,
      user: UserInfo.fromJson(data['user'] as Map<String, dynamic>),
    );
  }

  Future<bool> isAuthenticated() async {
    final token = await Storage.getToken();
    return token != null;
  }

  Future<void> logout() async {
    await Storage.clearAll();
  }
}

class AuthResult {
  final String token;
  final String tenantId;
  final UserInfo user;

  AuthResult({
    required this.token,
    required this.tenantId,
    required this.user,
  });
}

class UserInfo {
  final String id;
  final String email;
  final String name;
  final String role;

  UserInfo({
    required this.id,
    required this.email,
    required this.name,
    required this.role,
  });

  factory UserInfo.fromJson(Map<String, dynamic> json) => UserInfo(
        id: json['id'] as String,
        email: json['email'] as String,
        name: json['name'] as String,
        role: json['role'] as String,
      );
}

class OfflineService {
  final OperationRepo operationRepo;
  final PlotRepo plotRepo;
  final SyncQueueRepo syncQueueRepo;

  OfflineService({
    required this.operationRepo,
    required this.plotRepo,
    required this.syncQueueRepo,
  });

  Future<List<models.Operation>> getOperations() => operationRepo.getAll();
  Future<List<models.Plot>> getPlots() => plotRepo.getAll();

  Future<void> createOperation({
    required String plotId,
    required String type,
    required String date,
    String responsible = '',
    String productUsed = '',
    double quantity = 0,
    double cost = 0,
    String notes = '',
  }) async {
    final id = const Uuid().v4();
    final now = DateTime.now().toUtc().toIso8601String();

    final payload = {
      'plot_id': plotId,
      'type': type,
      'date': date,
      'responsible': responsible,
      'product_used': productUsed,
      'quantity': quantity,
      'cost': cost,
      'notes': notes,
    };

    await syncQueueRepo.enqueue('operation.created', payload);
    await operationRepo.insert(models.Operation(
      id: id,
      plotId: plotId,
      type: type,
      date: date,
      responsible: responsible,
      productUsed: productUsed,
      quantity: quantity,
      cost: cost,
      notes: notes,
      synced: 0,
      createdAt: now,
    ));
  }

  Future<void> refreshPlots() async {
    final response = await ApiClient.request('/plots');
    final data = response.data as List<dynamic>;
    final plots = data.map((j) => models.Plot.fromJson(j as Map<String, dynamic>)).toList();
    await plotRepo.upsertAll(plots);
  }
}

final authServiceProvider = Provider<AuthService>((ref) => AuthService());
final offlineServiceProvider = Provider<OfflineService>((ref) {
  final database = ref.read(db.databaseProvider);
  return OfflineService(
    operationRepo: OperationRepo(database),
    plotRepo: PlotRepo(database),
    syncQueueRepo: SyncQueueRepo(database),
  );
});
final syncServiceProvider = Provider<SyncService>((ref) {
  final database = ref.read(db.databaseProvider);
  return SyncService(SyncQueueRepo(database));
});
