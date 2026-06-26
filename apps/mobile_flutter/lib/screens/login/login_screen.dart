import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../services/services.dart';
import '../../shared/theme/app_theme.dart';

final loginControllerProvider =
    StateNotifierProvider<LoginController, LoginState>((ref) {
  return LoginController(ref.read(authServiceProvider));
});

class LoginState {
  final String email;
  final String password;
  final bool loading;
  final String? error;

  LoginState({
    this.email = 'joao@cafeos.com.br',
    this.password = '123456',
    this.loading = false,
    this.error,
  });

  LoginState copyWith({
    String? email,
    String? password,
    bool? loading,
    String? error,
  }) =>
      LoginState(
        email: email ?? this.email,
        password: password ?? this.password,
        loading: loading ?? this.loading,
        error: error,
      );
}

class LoginController extends StateNotifier<LoginState> {
  final AuthService _authService;

  LoginController(this._authService) : super(LoginState());

  void setEmail(String v) => state = state.copyWith(email: v, error: null);
  void setPassword(String v) => state = state.copyWith(password: v, error: null);

  Future<bool> login() async {
    if (state.email.isEmpty || state.password.isEmpty) {
      state = state.copyWith(error: 'Preencha todos os campos');
      return false;
    }

    state = state.copyWith(loading: true, error: null);

    try {
      await _authService.login(state.email, state.password);
      state = state.copyWith(loading: false);
      return true;
    } catch (e) {
      state = state.copyWith(loading: false, error: e.toString());
      return false;
    }
  }
}

class LoginScreen extends ConsumerWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(loginControllerProvider);
    final controller = ref.read(loginControllerProvider.notifier);

    return Scaffold(
      backgroundColor: AppTheme.cream,
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: Card(
              margin: EdgeInsets.zero,
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.coffee, size: 64, color: AppTheme.green),
                    const SizedBox(height: 8),
                    Text(
                      'CafeOS',
                      style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                            color: AppTheme.green,
                            fontWeight: FontWeight.bold,
                          ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      'App offline para operações de campo',
                      style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                            color: AppTheme.textSecondary,
                          ),
                    ),
                    const SizedBox(height: 32),
                    TextField(
                      decoration: const InputDecoration(
                        labelText: 'Email',
                        prefixIcon: Icon(Icons.email_outlined),
                      ),
                      keyboardType: TextInputType.emailAddress,
                      controller: TextEditingController(text: state.email),
                      onChanged: controller.setEmail,
                    ),
                    const SizedBox(height: 16),
                    TextField(
                      decoration: const InputDecoration(
                        labelText: 'Senha',
                        prefixIcon: Icon(Icons.lock_outlined),
                      ),
                      obscureText: true,
                      controller: TextEditingController(text: state.password),
                      onChanged: controller.setPassword,
                    ),
                    if (state.error != null) ...[
                      const SizedBox(height: 12),
                      Text(
                        state.error!,
                        style: const TextStyle(color: Colors.red),
                        textAlign: TextAlign.center,
                      ),
                    ],
                    const SizedBox(height: 24),
                    ElevatedButton(
                      onPressed: state.loading
                          ? null
                          : () async {
                              final ok = await controller.login();
                              if (ok && context.mounted) {
                                context.go('/operations');
                              }
                            },
                      child: state.loading
                          ? const SizedBox(
                              height: 20,
                              width: 20,
                              child: CircularProgressIndicator(
                                strokeWidth: 2,
                                color: Colors.white,
                              ),
                            )
                          : const Text('Entrar'),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
