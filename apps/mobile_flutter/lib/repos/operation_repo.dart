import 'package:uuid/uuid.dart';
import '../db/schema.dart';
import '../models/models.dart';

class OperationRepo {
  final AppDatabase _db;

  OperationRepo(this._db);

  Future<List<Operation>> getAll() async {
    final rows = await (_db.select(_db.operations)
          ..orderBy([(t) => OrderingTerm(expression: t.date, mode: OrderingMode.desc)]))
        .get();
    return rows.map(_fromRow).toList();
  }

  Future<void> insert(Operation op) async {
    await _db.into(_db.operations).insert(OperationsCompanion.insert(
          id: Value(op.id),
          plotId: Value(op.plotId),
          type: Value(op.type),
          date: Value(op.date),
          responsible: Value(op.responsible),
          productUsed: Value(op.productUsed),
          quantity: Value(op.quantity),
          cost: Value(op.cost),
          notes: Value(op.notes),
        ));
  }

  Future<void> markSynced(String id) async {
    await (_db.update(_db.operations)
          ..where((t) => t.id.equals(id)))
        .write(const OperationsCompanion(synced: Value(1)));
  }

  Operation _fromRow(OperationsData row) => Operation(
        id: row.id,
        plotId: row.plotId,
        type: row.type,
        date: row.date,
        responsible: row.responsible,
        productUsed: row.productUsed,
        quantity: row.quantity,
        cost: row.cost,
        notes: row.notes,
        synced: row.synced,
        createdAt: row.createdAt,
      );
}
