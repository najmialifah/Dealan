import React, { useEffect, useState } from 'react';
import { View, Text, StyleSheet, ActivityIndicator, Alert, Button } from 'react-native';
import { findDriver } from '../services/matchingApi';

export default function MatchingScreen({ route, navigation }) {
  const { order_id } = route.params || {};
  const [status, setStatus] = useState('Mencari driver di sekitar...');

  useEffect(() => {
    const matchDriver = async () => {
      try {
        const payload = {
          order_id: String(order_id),
          lat: -6.2088,
          lon: 106.8456,
          radius: 5000,
          service_type: 'ride'
        };
        const res = await findDriver(payload);
        setStatus(`Driver Ditemukan: ${res.driver_id || res.message}`);
        // Wait a bit, then move to Payment
        setTimeout(() => {
          navigation.navigate('Payment', { order_id, driver_id: res.driver_id || 'drv-123', nominal: 20000 });
        }, 2000);
      } catch (err) {
        setStatus('Gagal mencari driver.');
      }
    };

    matchDriver();
  }, []);

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Mencari Driver</Text>
      <ActivityIndicator size="large" color="#00ff00" style={{ marginVertical: 20 }} />
      <Text style={styles.status}>{status}</Text>
      
      {status.includes('Gagal') && (
        <Button title="Lewati ke Pembayaran (Demo)" onPress={() => navigation.navigate('Payment', { order_id, driver_id: 'drv-demo', nominal: 20000 })} />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 20, justifyContent: 'center', alignItems: 'center' },
  title: { fontSize: 22, fontWeight: 'bold', marginBottom: 10 },
  status: { fontSize: 16, textAlign: 'center', marginBottom: 20 }
});
