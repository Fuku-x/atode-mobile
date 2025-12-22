import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/material.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();

  bool _isSubmitting = false;
  String? _errorMessage;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _submit({required bool createAccount}) async {
    if (_isSubmitting) return;

    setState(() {
      _errorMessage = null;
    });

    final isValid = _formKey.currentState?.validate() ?? false;
    if (!isValid) return;

    setState(() {
      _isSubmitting = true;
    });

    try {
      final email = _emailController.text.trim();
      final password = _passwordController.text;

      if (createAccount) {
        await FirebaseAuth.instance.createUserWithEmailAndPassword(
          email: email,
          password: password,
        );
      } else {
        await FirebaseAuth.instance.signInWithEmailAndPassword(
          email: email,
          password: password,
        );
      }
    } on FirebaseAuthException catch (e) {
      setState(() {
        _errorMessage = _firebaseMessage(e);
      });
    } catch (_) {
      setState(() {
        _errorMessage = 'ログインに失敗しました';
      });
    } finally {
      if (mounted) {
        setState(() {
          _isSubmitting = false;
        });
      }
    }
  }

  String _firebaseMessage(FirebaseAuthException e) {
    switch (e.code) {
      case 'invalid-email':
        return 'メールアドレスの形式が正しくありません';
      case 'user-disabled':
        return 'このユーザーは無効化されています';
      case 'user-not-found':
      case 'wrong-password':
      case 'invalid-credential':
        return 'メールアドレスまたはパスワードが違います';
      case 'email-already-in-use':
        return 'このメールアドレスは既に使われています';
      case 'weak-password':
        return 'パスワードが弱すぎます';
      default:
        return '認証に失敗しました (${e.code})';
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Login')),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextFormField(
                  controller: _emailController,
                  keyboardType: TextInputType.emailAddress,
                  autofillHints: const [
                    AutofillHints.username,
                    AutofillHints.email,
                  ],
                  decoration: const InputDecoration(
                    labelText: 'メールアドレス',
                  ),
                  validator: (value) {
                    final v = (value ?? '').trim();
                    if (v.isEmpty) return 'メールアドレスを入力してください';
                    if (!v.contains('@')) return 'メールアドレスの形式が正しくありません';
                    return null;
                  },
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: _passwordController,
                  obscureText: true,
                  autofillHints: const [AutofillHints.password],
                  decoration: const InputDecoration(
                    labelText: 'パスワード',
                  ),
                  validator: (value) {
                    final v = value ?? '';
                    if (v.isEmpty) return 'パスワードを入力してください';
                    if (v.length < 6) return '6文字以上で入力してください';
                    return null;
                  },
                ),
                const SizedBox(height: 12),
                if (_errorMessage != null) ...[
                  Text(
                    _errorMessage!,
                    style: TextStyle(color: Theme.of(context).colorScheme.error),
                  ),
                  const SizedBox(height: 12),
                ],
                ElevatedButton(
                  onPressed:
                      _isSubmitting ? null : () => _submit(createAccount: false),
                  child: _isSubmitting
                      ? const SizedBox(
                          width: 18,
                          height: 18,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Text('ログイン'),
                ),
                const SizedBox(height: 8),
                OutlinedButton(
                  onPressed:
                      _isSubmitting ? null : () => _submit(createAccount: true),
                  child: const Text('新規登録'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
