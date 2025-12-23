import 'package:firebase_auth/firebase_auth.dart';

String firebaseAuthErrorMessage(FirebaseAuthException e) {
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
