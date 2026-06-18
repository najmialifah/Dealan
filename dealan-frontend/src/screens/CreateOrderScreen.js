import React, { useState } from 'react';
import { View, Text, TextInput, Button, StyleSheet, ActivityIndicator, Alert, ScrollView } from 'react-native';
import { createOrder } from '../services/orderApi';

export default function CreateOrderScreen({ navigation }) {
  const [serviceType, setServiceType] = useState('ride'); // ride, car, send
  const [userId, setUserId] = useState('1'); // Placeholder for User ID
  const [detailPaket, setDetailPaket] = useState('');
  const [loading, setLoading] = useState(false);

  const handleCreateOrder = async () => {
    try {
      setLoading(true);
      const payload = {
        user_id: parseInt(userId),
        status: 'PENDING',
        detail_paket: {
          service_type: serviceType,
          description: serviceType === 'send' ? detailPaket : '',
        }
      };

      const res = await createOrder(payload);
      Alert.alert('Success', 'Pesanan berhasil dibuat', [
        { text: 'OK', onPress: () => navigation.navigate('Negotiation', { order_id: res.data.id }) }
      ]);
    } catch (err) {
      // Handled globally
    } finally {
      setLoading(false);
    }
  };

  return (
    <ScrollView contentContainerStyle={styles.container}>
      <Text style={styles.title}>Buat Pesanan Baru</Text>

      <View style={styles.serviceContainer}>
        <Button title="GoRide" color={serviceType === 'ride' ? 'green' : 'gray'} onPress={() => setServiceType('ride')} />
        <Button title="GoCar" color={serviceType === 'car' ? 'green' : 'gray'} onPress={() => setServiceType('car')} />
        <Button title="GoSend" color={serviceType === 'send' ? 'green' : 'gray'} onPress={() => setServiceType('send')} />
      </View>

      {serviceType === 'send' && (
        <TextInput
          style={styles.input}
          placeholder="Detail barang (contoh: Dokumen penting)"
          value={detailPaket}
          onChangeText={setDetailPaket}
        />
      )}

      {loading ? (
        <ActivityIndicator size="large" color="#0000ff" />
      ) : (
        <Button title="Pesan Sekarang" onPress={handleCreateOrder} />
      )}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: { flexGrow: 1, padding: 20 },
  title: { fontSize: 22, fontWeight: 'bold', marginBottom: 20, textAlign: 'center' },
  serviceContainer: { flexDirection: 'row', justifyContent: 'space-between', marginBottom: 20 },
  input: { borderWidth: 1, borderColor: '#ccc', padding: 10, marginBottom: 15, borderRadius: 5 }
});
