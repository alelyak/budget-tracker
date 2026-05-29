import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, TextInput, Button, Title, Text, Paper, Box } from '@mantine/core';
import { useForm } from '@mantine/form';
import api from '../api/axios';

function Login() {
    const navigate = useNavigate();
    const [isLogin, setIsLogin] = useState(true);
    const [error, setError] = useState('');

    const form = useForm({
        initialValues: { email: '', password: '' },
        validate: {
            email: (val) => (/^\S+@\S+$/.test(val) ? null : 'Некорректный email'),
            password: (val) => (val.length < 4 ? 'Пароль должен быть длиннее 4 символов' : null),
        },
    });

    const handleSubmit = async ({ email, password }) => {
        setError('');
        try {
            const endpoint = isLogin ? '/login' : '/register';
            const res = await api.post(endpoint, { email, password });
            localStorage.setItem('token', res.data.token);
            localStorage.setItem('user', JSON.stringify(res.data.user));
            navigate('/dashboard');
        } catch (err) {
            setError(err.response?.data?.error || 'Ошибка соединения');
        }
    };

    return (
        <Box style={{ height: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#f5f5f5' }}>
            <Card shadow="lg" padding="xl" radius="md" withBorder style={{ width: 420 }}>
                <Title order={2} ta="center" mb="md">💰 Финансы</Title>
                <Title order={4} ta="center" c="dimmed" mb="xl">
                    {isLogin ? 'Войдите в аккаунт' : 'Создайте аккаунт'}
                </Title>

                <form onSubmit={form.onSubmit(handleSubmit)}>
                    <TextInput
                        label="Email"
                        placeholder="user@example.com"
                        {...form.getInputProps('email')}
                        mb="sm"
                    />
                    <TextInput
                        label="Пароль"
                        type="password"
                        {...form.getInputProps('password')}
                        mb="md"
                    />

                    {error && (
                        <Paper bg="red.1" c="red" p="xs" mb="md" radius="sm">
                            <Text size="sm">{error}</Text>
                        </Paper>
                    )}

                    <Button type="submit" fullWidth color="blue">
                        {isLogin ? 'Войти' : 'Зарегистрироваться'}
                    </Button>

                    <Button
                        variant="subtle"
                        fullWidth
                        mt="md"
                        onClick={() => setIsLogin(!isLogin)}
                    >
                        {isLogin ? 'Нет аккаунта? Зарегистрироваться' : 'Уже есть аккаунт? Войти'}
                    </Button>
                </form>
            </Card>
        </Box>
    );
}

export default Login;
