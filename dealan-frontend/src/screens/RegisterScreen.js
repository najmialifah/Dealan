import React, { useState } from 'react';
import { View, Text, TextInput, Button, StyleSheet, Alert, ActivityIndicator, ScrollView } from 'react-native';
import { register } from '../services/authApi';

export default function RegisterScreen({ navigation }) {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('user'); // default role
  const [loading, setLoading] = useState(false);

  const handleRegister = async () => {
    if (!name || !email || !password) {
      Alert.alert('Error', 'Semua field harus diisi');
      return;
    }
    try {
      setLoading(true);
      const res = await register({ name, email, password, role });
      Alert.alert('Success', 'Registrasi berhasil, silakan login', [
        { text: 'OK', onPress: () => navigation.goBack() }
      ]);
    } catch (err) {
      // Error is caught globally
    } finally {
      setLoading(false);
    }
  };

  return (
    <ScrollView contentContainerStyle={styles.container}>
      <Text style={styles.title}>Daftar Akun Baru</Text>
      <TextInput
        style={styles.input}
        placeholder="Nama Lengkap"
        value={name}
        onChangeText={setName}
      />
      <TextInput
        style={styles.input}
        placeholder="Email"
        value={email}
        onChangeText={setEmail}
        autoCapitalize="none"
        keyboardType="email-address"
      />
      <TextInput
        style={styles.input}
        placeholder="Password"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
      />
      {/* Simplify role selection for demo */}
      <View style={styles.roleContainer}>
        <Button title="Role: User" color={role === 'user' ? 'blue' : 'gray'} onPress={() => setRole('user')} />
        <Button title="Role: Driver" color={role === 'driver' ? 'blue' : 'gray'} onPress={() => setRole('driver')} />
      </View>
      
      {loading ? (
        <ActivityIndicator size="large" color="#0000ff" />
      ) : (
        <Button title="Daftar" onPress={handleRegister} />
      )}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: { flexGrow: 1, padding: 20, justifyContent: 'center' },
  title: { fontSize: 24, fontWeight: 'bold', marginBottom: 20, textAlign: 'center' },
  input: { borderWidth: 1, borderColor: '#ccc', padding: 10, marginBottom: 15, borderRadius: 5 },
  roleContainer: { flexDirection: 'row', justifyContent: 'space-around', marginBottom: 20 }
});
