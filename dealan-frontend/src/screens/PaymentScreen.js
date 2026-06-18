import React, { useState } from 'react';
import { View, Text, Button, StyleSheet, ActivityIndicator, Alert } from 'react-native';
import { processPayment } from '../services/paymentApi';

export default function PaymentScreen({ route, navigation }) {
  const { order_id, driver_id, nominal } = route.params || {};
  const [method, setMethod] = useState('cash'); // cash, qris, bank
  const [loading, setLoading] = useState(false);

  const handlePayment = async () => {
    try {
      setLoading(true);
      const payload = {
        order_id: String(order_id),
        nominal: parseFloat(nominal),
        metode_pembayaran: method,
        user_id: 'usr-1', // Mock user id
        driver_id: String(driver_id)
      };

      const res = await processPayment(payload);
      Alert.alert('Success', 'Pembayaran Berhasil / Terproses', [
        { text: 'Lanjut ke Rating', onPress: () => navigation.navigate('Rating', { order_id, driver_id }) }
      ]);
    } catch (err) {
      // Handled globally
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Pembayaran</Text>
      <Text style={styles.info}>Total Tagihan: Rp {nominal}</Text>
      
      <Text style={styles.subtitle}>Pilih Metode Pembayaran:</Text>
      <View style={styles.methodContainer}>
        <Button title="Tunai (Cash)" color={method === 'cash' ? 'blue' : 'gray'} onPress={() => setMethod('cash')} />
        <Button title="QRIS" color={method === 'qris' ? 'blue' : 'gray'} onPress={() => setMethod('qris')} />
        <Button title="Transfer Bank" color={method === 'bank' ? 'blue' : 'gray'} onPress={() => setMethod('bank')} />
      </View>

      {loading ? (
        <ActivityIndicator size="large" color="#0000ff" />
      ) : (
        <Button title="Bayar Sekarang" onPress={handlePayment} />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 20, justifyContent: 'center' },
  title: { fontSize: 24, fontWeight: 'bold', marginBottom: 10, textAlign: 'center' },
  info: { fontSize: 18, marginBottom: 20, textAlign: 'center' },
  subtitle: { fontSize: 16, marginBottom: 10 },
  methodContainer: { flexDirection: 'row', justifyContent: 'space-between', marginBottom: 30 }
});
