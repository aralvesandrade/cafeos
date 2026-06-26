import 'package:drift/drift.dart';
import '../db/schema.dart' as db;
import '../models/models.dart' as models;

class OperationRepo {
  final db.AppDatabase _db;

  OperationRepo(this._db);

  Future<List<models.Operation>> getAll() async {
    final rows = await (_db.select(_db.operations)
          ..orderBy([(t) => OrderingTerm(expression: t.date, mode: OrderingMode.desc)]))
        .get();
    return rows.map(_fromRow).toList();
  }

  Future<void> insert(models.Operation op) async {
    await _db.into(_db.operations).insert(db.OperationsCompanion.insert(
          id: op.id,
          plotId: op.plotId,
          type: op.type,
          date: op.date,
          responsible: Value(op.responsible),
          productUsed: Value(op.productUsed),
          quantity: Value(op.quantity),
          cost: Value(op.cost),
          notes: Value(op.notes),
          synced: Value(op.synced),
          createdAt: op.createdAt,
        ));
  }

  Future<void> markSynced(String id) async {
    await (_db.update(_db.operations)
          ..where((t) => t.id.equals(id)))
        .write(const db.OperationsCompanion(synced: Value(1)));
  }

  models.Operation _fromRow(db.Operation row) => models.Operation(
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
