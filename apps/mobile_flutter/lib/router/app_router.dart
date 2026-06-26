import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../screens/home/home_screen.dart';
import '../screens/login/login_screen.dart';
import '../screens/operations/operations_screen.dart';
import '../screens/pending_sync/pending_sync_screen.dart';

final GlobalKey<NavigatorState> _rootNavigator = GlobalKey<NavigatorState>();

final appRouter = GoRouter(
  navigatorKey: _rootNavigator,
  initialLocation: '/login',
  routes: [
    GoRoute(
      path: '/login',
      builder: (context, state) => const LoginScreen(),
    ),
    ShellRoute(
      builder: (context, state, child) => HomeScreen(child: child),
      routes: [
        GoRoute(
          path: '/',
          redirect: (context, state) => '/operations',
        ),
        GoRoute(
          path: '/operations',
          pageBuilder: (context, state) => NoTransitionPage(
            key: state.pageKey,
            child: const OperationsScreen(),
          ),
        ),
        GoRoute(
          path: '/pending-sync',
          pageBuilder: (context, state) => NoTransitionPage(
            key: state.pageKey,
            child: const PendingSyncScreen(),
          ),
        ),
      ],
    ),
  ],
);
