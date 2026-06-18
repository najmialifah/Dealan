import React, { useState } from 'react';
import { View, Text, TextInput, Button, StyleSheet, ActivityIndicator, Alert } from 'react-native';
import { negotiatePrice } from '../services/pricingApi';

export default function NegotiationScreen({ route, navigation }) {
  const { order_id } = route.params || {};
  const [requestedPrice, setRequestedPrice] = useState('');
  const [loading, setLoading] = useState(false);
  const originalPrice = 20000; // Mock estimate for demo

  const handleNegotiate = async () => {
    if (!requestedPrice) {
      Alert.alert('Error', 'Masukkan harga tawaran');
      return;
    }
    try {
      setLoading(true);
      const payload = {
        order_id: String(order_id),
        original_price: originalPrice,
        requested_price: parseFloat(requestedPrice)
      };

      const res = await negotiatePrice(payload);
      if (res.status === 'approved' || res.status === 'accepted') {
        Alert.alert('Success', 'Tawaran Diterima!', [
          { text: 'OK', onPress: () => navigation.navigate('Matching', { order_id }) }
        ]);
      } else {
        Alert.alert('Ditolak', 'Harga terlalu rendah, coba naikkan sedikit.');
      }
    } catch (err) {
      // Globally handled
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Negosiasi Harga</Text>
      <Text style={styles.info}>ID Pesanan: {order_id}</Text>
      <Text style={styles.info}>Estimasi Harga: Rp {originalPrice}</Text>

      <TextInput
        style={styles.input}
        placeholder="Tawarkan harga Anda"
        keyboardType="numeric"
        value={requestedPrice}
        onChangeText={setRequestedPrice}
      />

      {loading ? (
        <ActivityIndicator size="large" color="#0000ff" />
      ) : (
        <Button title="Ajukan Tawaran" onPress={handleNegotiate} />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 20, justifyContent: 'center' },
  title: { fontSize: 22, fontWeight: 'bold', marginBottom: 20, textAlign: 'center' },
  info: { fontSize: 16, marginBottom: 10 },
  input: { borderWidth: 1, borderColor: '#ccc', padding: 10, marginBottom: 20, borderRadius: 5 }
});
