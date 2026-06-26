import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../services/services.dart';
import '../../shared/theme/app_theme.dart';
import '../../models/models.dart';

final operationsControllerProvider =
    StateNotifierProvider<OperationsController, OperationsState>((ref) {
  return OperationsController(ref.read(offlineServiceProvider));
});

class OperationsState {
  final List<Operation> operations;
  final List<Plot> plots;
  final bool loading;
  final String? error;

  OperationsState({
    this.operations = const [],
    this.plots = const [],
    this.loading = true,
    this.error,
  });

  OperationsState copyWith({
    List<Operation>? operations,
    List<Plot>? plots,
    bool? loading,
    String? error,
  }) =>
      OperationsState(
        operations: operations ?? this.operations,
        plots: plots ?? this.plots,
        loading: loading ?? this.loading,
        error: error,
      );
}

class OperationsController extends StateNotifier<OperationsState> {
  final OfflineService _service;

  OperationsController(this._service) : super(OperationsState()) {
    load();
  }

  Future<void> load() async {
    state = state.copyWith(loading: true, error: null);
    try {
      final ops = await _service.getOperations();
      final plots = await _service.getPlots();
      state = state.copyWith(operations: ops, plots: plots, loading: false);
    } catch (e) {
      state = state.copyWith(loading: false, error: e.toString());
    }
  }

  Future<void> refreshPlots() async {
    try {
      await _service.refreshPlots();
      final plots = await _service.getPlots();
      state = state.copyWith(plots: plots);
    } catch (_) {}
  }

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
    await _service.createOperation(
      plotId: plotId,
      type: type,
      date: date,
      responsible: responsible,
      productUsed: productUsed,
      quantity: quantity,
      cost: cost,
      notes: notes,
    );
    await load();
  }
}

class OperationsScreen extends ConsumerStatefulWidget {
  const OperationsScreen({super.key});

  @override
  ConsumerState<OperationsScreen> createState() => _OperationsScreenState();
}

class _OperationsScreenState extends ConsumerState<OperationsScreen> {
  @override
  void initState() {
    super.initState();
    ref.read(operationsControllerProvider.notifier).refreshPlots();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(operationsControllerProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Operações')),
      body: _buildBody(state),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showForm(context),
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildBody(OperationsState state) {
    if (state.loading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (state.error != null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(state.error!, style: const TextStyle(color: Colors.red)),
            const SizedBox(height: 16),
            ElevatedButton(onPressed: () => ref.read(operationsControllerProvider.notifier).load(), child: const Text('Tentar novamente')),
          ],
        ),
      );
    }

    if (state.operations.isEmpty) {
      return const Center(
        child: Text(
          'Nenhuma operação registrada',
          style: TextStyle(color: AppTheme.textSecondary, fontSize: 16),
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: () => ref.read(operationsControllerProvider.notifier).load(),
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: state.operations.length,
        itemBuilder: (context, index) {
          final op = state.operations[index];
          return _OperationCard(operation: op);
        },
      ),
    );
  }

  void _showForm(BuildContext context) {
    final state = ref.read(operationsControllerProvider);
    final controller = ref.read(operationsControllerProvider.notifier);

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) => _OperationForm(
        plots: state.plots,
        onSave: (data) {
          controller.createOperation(
            plotId: data['plotId']!,
            type: data['type']!,
            date: data['date']!,
            responsible: data['responsible'] ?? '',
            productUsed: data['productUsed'] ?? '',
            quantity: double.tryParse(data['quantity'] ?? '0') ?? 0,
            cost: double.tryParse(data['cost'] ?? '0') ?? 0,
            notes: data['notes'] ?? '',
          );
          Navigator.of(ctx).pop();
        },
      ),
    );
  }
}

class _OperationCard extends StatelessWidget {
  final Operation operation;

  const _OperationCard({required this.operation});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              width: 4,
              height: 48,
              decoration: BoxDecoration(
                color: _typeColor(operation.type),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: _typeColor(operation.type).withValues(alpha: 0.15),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text(
                          _typeLabel(operation.type),
                          style: TextStyle(
                            fontSize: 12,
                            fontWeight: FontWeight.w600,
                            color: _typeColor(operation.type),
                          ),
                        ),
                      ),
                      if (operation.synced == 0) ...[
                        const SizedBox(width: 8),
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                          decoration: BoxDecoration(
                            color: Colors.orange.withValues(alpha: 0.15),
                            borderRadius: BorderRadius.circular(4),
                          ),
                          child: const Text(
                            'Pendente',
                            style: TextStyle(fontSize: 11, color: Colors.orange),
                          ),
                        ),
                      ],
                    ],
                  ),
                  const SizedBox(height: 4),
                  Text(
                    operation.date,
                    style: const TextStyle(
                      fontSize: 13,
                      color: AppTheme.textSecondary,
                    ),
                  ),
                  if (operation.responsible.isNotEmpty)
                    Text(
                      operation.responsible,
                      style: const TextStyle(fontSize: 13),
                    ),
                  if (operation.cost > 0)
                    Text(
                      'R\$ ${operation.cost.toStringAsFixed(2)}',
                      style: const TextStyle(
                        fontSize: 13,
                        fontWeight: FontWeight.w600,
                        color: AppTheme.green,
                      ),
                    ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Color _typeColor(String type) {
    switch (type) {
      case 'adubacao':
        return Colors.blue;
      case 'pulverizacao':
        return Colors.orange;
      case 'irrigacao':
        return Colors.green;
      case 'poda':
        return Colors.brown;
      case 'colheita':
        return Colors.red;
      default:
        return Colors.grey;
    }
  }

  String _typeLabel(String type) {
    switch (type) {
      case 'adubacao':
        return 'Adubação';
      case 'pulverizacao':
        return 'Pulverização';
      case 'irrigacao':
        return 'Irrigação';
      case 'poda':
        return 'Poda';
      case 'colheita':
        return 'Colheita';
      default:
        return type;
    }
  }
}

class _OperationForm extends StatefulWidget {
  final List<Plot> plots;
  final void Function(Map<String, String> data) onSave;

  const _OperationForm({required this.plots, required this.onSave});

  @override
  State<_OperationForm> createState() => _OperationFormState();
}

class _OperationFormState extends State<_OperationForm> {
  final _formKey = GlobalKey<FormState>();
  String _type = 'adubacao';
  String _plotId = '';
  final _dateCtrl = TextEditingController(text: DateTime.now().toIso8601String().substring(0, 10));
  final _responsibleCtrl = TextEditingController();
  final _productCtrl = TextEditingController();
  final _qtyCtrl = TextEditingController();
  final _costCtrl = TextEditingController();
  final _notesCtrl = TextEditingController();

  static const _types = [
    {'value': 'adubacao', 'label': 'Adubação'},
    {'value': 'pulverizacao', 'label': 'Pulverização'},
    {'value': 'irrigacao', 'label': 'Irrigação'},
    {'value': 'poda', 'label': 'Poda'},
    {'value': 'colheita', 'label': 'Colheita'},
  ];

  @override
  void dispose() {
    _dateCtrl.dispose();
    _responsibleCtrl.dispose();
    _productCtrl.dispose();
    _qtyCtrl.dispose();
    _costCtrl.dispose();
    _notesCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.only(
        bottom: MediaQuery.of(context).viewInsets.bottom,
        left: 24,
        right: 24,
        top: 24,
      ),
      child: Form(
        key: _formKey,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text('Nova Operação', style: Theme.of(context).textTheme.titleLarge),
              const SizedBox(height: 16),
              DropdownButtonFormField<String>(
                value: _plotId.isEmpty ? null : _plotId,
                decoration: const InputDecoration(labelText: 'Talhão'),
                items: widget.plots
                    .map((p) => DropdownMenuItem(value: p.id, child: Text(p.name)))
                    .toList(),
                onChanged: (v) => _plotId = v ?? '',
                validator: (v) => v == null || v.isEmpty ? 'Selecione um talhão' : null,
              ),
              const SizedBox(height: 16),
              const Text('Tipo', style: TextStyle(fontWeight: FontWeight.w500)),
              const SizedBox(height: 8),
              Wrap(
                spacing: 8,
                children: _types.map((t) => ChoiceChip(
                      label: Text(t['label']!),
                      selected: _type == t['value'],
                      onSelected: (_) => setState(() => _type = t['value']!),
                    )).toList(),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _dateCtrl,
                decoration: const InputDecoration(labelText: 'Data (YYYY-MM-DD)'),
                validator: (v) => v == null || v.isEmpty ? 'Campo obrigatório' : null,
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _responsibleCtrl,
                decoration: const InputDecoration(labelText: 'Responsável'),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _productCtrl,
                decoration: const InputDecoration(labelText: 'Produto utilizado'),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _qtyCtrl,
                decoration: const InputDecoration(labelText: 'Quantidade'),
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _costCtrl,
                decoration: const InputDecoration(labelText: 'Custo (R\$)'),
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _notesCtrl,
                decoration: const InputDecoration(labelText: 'Observações'),
                maxLines: 3,
              ),
              const SizedBox(height: 24),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton(
                      onPressed: () => Navigator.of(context).pop(),
                      child: const Text('Cancelar'),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: ElevatedButton(
                      onPressed: () {
                        if (_formKey.currentState!.validate()) {
                          widget.onSave({
                            'plotId': _plotId,
                            'type': _type,
                            'date': _dateCtrl.text,
                            'responsible': _responsibleCtrl.text,
                            'productUsed': _productCtrl.text,
                            'quantity': _qtyCtrl.text,
                            'cost': _costCtrl.text,
                            'notes': _notesCtrl.text,
                          });
                        }
                      },
                      child: const Text('Salvar'),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }
}
