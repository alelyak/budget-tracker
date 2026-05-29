import { useState, useEffect } from "react";
import {
  Title,
  Text,
  Card,
  Table,
  Button,
  Modal,
  TextInput,
  Select,
  Group,
  Loader,
  Center,
  Stack,
  NumberInput,
  Textarea,
  Alert,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import {
  IconAlertCircle,
  IconPlus,
  IconEdit,
  IconTrash,
  IconCheck,
  IconX,
} from "@tabler/icons-react";
import api from "../api/axios";

function Transactions() {
  const [transactions, setTransactions] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [modalOpened, setModalOpened] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [filters, setFilters] = useState({
    start_date: "",
    end_date: "",
    category_id: "",
  });

  const form = useForm({
    initialValues: {
      category_id: "",
      amount: "",
      date: "",
      note: "",
    },
    validate: {
      category_id: (val) => (!val ? "Выберите категорию" : null),
      amount: (val) => (!val || val <= 0 ? "Сумма должна быть больше 0" : null),
      date: (val) => (!val ? "Выберите дату" : null),
    },
  });

  const fetchTransactions = async () => {
    setLoading(true);
    setError(null);
    try {
      const params = {};
      if (filters.start_date) params.start_date = filters.start_date;
      if (filters.end_date) params.end_date = filters.end_date;
      if (filters.category_id) params.category_id = filters.category_id;

      const res = await api.get("/transactions", { params });
      setTransactions(res.data || []);
    } catch (error) {
      setError(
        "Не удалось загрузить транзакции. Проверьте подключение к серверу.",
      );
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const res = await api.get("/categories");
      setCategories(res.data || []);
    } catch (error) {
      console.error("Error fetching categories:", error);
    }
  };

  useEffect(() => {
    fetchTransactions();
    fetchCategories();
  }, [filters]);

  const handleSubmit = async (values) => {
    try {
      const data = {
        category_id: parseInt(values.category_id),
        amount: parseFloat(values.amount),
        date: values.date,
        note: values.note || "",
      };

      if (editingId) {
        await api.put(`/transactions/${editingId}`, data);
        notifications.show({
          title: "Успешно",
          message: "Транзакция обновлена",
          color: "green",
          icon: <IconCheck size={16} />,
          autoClose: 3000,
        });
      } else {
        const response = await api.post("/transactions", data);
        if (response.data.budget_warning) {
          notifications.show({
            title: "Внимание",
            message: response.data.budget_warning.message,
            color: "yellow",
            icon: <IconAlertCircle size={16} />,
            autoClose: 5000,
          });
        } else {
          notifications.show({
            title: "Успешно",
            message: "Транзакция создана",
            color: "green",
            icon: <IconCheck size={16} />,
            autoClose: 3000,
          });
        }
      }
      closeModal();
      fetchTransactions();
    } catch (error) {
      notifications.show({
        title: "Ошибка",
        message: error.response?.data?.error || "Ошибка при сохранении",
        color: "red",
        icon: <IconX size={16} />,
        autoClose: 4000,
      });
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm("Удалить транзакцию?")) {
      try {
        await api.delete(`/transactions/${id}`);
        notifications.show({
          title: "Успешно",
          message: "Транзакция удалена",
          color: "green",
          icon: <IconCheck size={16} />,
          autoClose: 3000,
        });
        fetchTransactions();
      } catch (error) {
        notifications.show({
          title: "Ошибка",
          message: "Не удалось удалить транзакцию",
          color: "red",
          icon: <IconX size={16} />,
          autoClose: 4000,
        });
      }
    }
  };

  const openModal = (transaction = null) => {
    if (transaction) {
      setEditingId(transaction.id);
      form.setValues({
        category_id: transaction.category_id.toString(),
        amount: transaction.amount,
        date: transaction.date.split("T")[0],
        note: transaction.note || "",
      });
    } else {
      setEditingId(null);
      form.reset();
      form.setValues({ date: new Date().toISOString().split("T")[0] });
    }
    setModalOpened(true);
  };

  const closeModal = () => {
    setModalOpened(false);
    form.reset();
    setEditingId(null);
  };

  if (loading) {
    return (
      <Center style={{ height: "60vh", flexDirection: "column" }}>
        <Loader size="xl" mb="md" />
        <Text c="dimmed">Загрузка транзакций...</Text>
      </Center>
    );
  }

  if (error) {
    return (
      <Center style={{ height: "60vh" }}>
        <Alert
          icon={<IconAlertCircle size={16} />}
          title="Ошибка"
          color="red"
          variant="filled"
        >
          {error}
        </Alert>
      </Center>
    );
  }

  return (
    <>
      <Group justify="space-between" mb="md">
        <Title order={2}>💰 Транзакции</Title>
        <Button
          leftSection={<IconPlus size={20} />}
          onClick={() => openModal()}
        >
          Добавить
        </Button>
      </Group>

      <Card shadow="sm" padding="md" radius="md" withBorder mb="md">
        <Group grow align="flex-end">
          <TextInput
            label="Дата от"
            type="date"
            value={filters.start_date}
            onChange={(e) =>
              setFilters({ ...filters, start_date: e.target.value })
            }
          />
          <TextInput
            label="Дата до"
            type="date"
            value={filters.end_date}
            onChange={(e) =>
              setFilters({ ...filters, end_date: e.target.value })
            }
          />
          <Select
            label="Категория"
            placeholder="Все"
            data={[
              { value: "", label: "Все категории" },
              ...categories.map((c) => ({
                value: c.id.toString(),
                label: c.name,
              })),
            ]}
            value={filters.category_id}
            onChange={(value) => setFilters({ ...filters, category_id: value })}
          />
          <Button
            variant="subtle"
            onClick={() =>
              setFilters({ start_date: "", end_date: "", category_id: "" })
            }
          >
            Сбросить
          </Button>
        </Group>
      </Card>

      <Card shadow="sm" padding="md" radius="md" withBorder>
        {transactions.length === 0 ? (
          <Center style={{ flexDirection: "column", padding: "60px 20px" }}>
            <IconAlertCircle size={48} color="#868e96" stroke={1.5} />
            <Title order={3} mt="md" c="dimmed">
              Нет транзакций
            </Title>
            <Text c="dimmed" ta="center" mt="sm">
              У вас пока нет ни одной транзакции
            </Text>
            <Button
              leftSection={<IconPlus size={20} />}
              onClick={() => openModal()}
              variant="outline"
              mt="md"
            >
              Добавить первую транзакцию
            </Button>
          </Center>
        ) : (
          <Table striped highlightOnHover>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Дата</Table.Th>
                <Table.Th>Категория</Table.Th>
                <Table.Th>Сумма</Table.Th>
                <Table.Th>Примечание</Table.Th>
                <Table.Th style={{ width: 80 }}>Действия</Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {transactions.map((tx) => (
                <Table.Tr key={tx.id}>
                  <Table.Td>{tx.date?.split("T")[0]}</Table.Td>
                  <Table.Td>{tx.category_name}</Table.Td>
                  <Table.Td
                    style={{
                      color: tx.category_type === "income" ? "green" : "red",
                      fontWeight: 600,
                    }}
                  >
                    {tx.category_type === "income" ? "+" : "-"}
                    {tx.amount.toLocaleString()} ₽
                  </Table.Td>
                  <Table.Td>{tx.note || "—"}</Table.Td>
                  <Table.Td>
                    <Group gap="xs">
                      <Button
                        size="xs"
                        variant="subtle"
                        color="blue"
                        onClick={() => openModal(tx)}
                      >
                        <IconEdit size={16} />
                      </Button>
                      <Button
                        size="xs"
                        variant="subtle"
                        color="red"
                        onClick={() => handleDelete(tx.id)}
                      >
                        <IconTrash size={16} />
                      </Button>
                    </Group>
                  </Table.Td>
                </Table.Tr>
              ))}
            </Table.Tbody>
          </Table>
        )}
      </Card>

      <Modal
        opened={modalOpened}
        onClose={closeModal}
        title={editingId ? "Редактировать транзакцию" : "Новая транзакция"}
        size="md"
      >
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <Stack>
            <Select
              label="Категория"
              placeholder="Выберите категорию"
              data={categories.map((c) => ({
                value: c.id.toString(),
                label: `${c.type === "income" ? "💰" : "💸"} ${c.name}`,
              }))}
              required
              {...form.getInputProps("category_id")}
            />
            <NumberInput
              label="Сумма (₽)"
              placeholder="Введите сумму"
              min={1}
              required
              {...form.getInputProps("amount")}
            />
            <TextInput
              label="Дата"
              type="date"
              required
              {...form.getInputProps("date")}
            />
            <Textarea
              label="Примечание"
              placeholder="Комментарий (необязательно)"
              {...form.getInputProps("note")}
            />
            <Group justify="flex-end" mt="md">
              <Button variant="subtle" onClick={closeModal}>
                Отмена
              </Button>
              <Button type="submit" color="blue">
                {editingId ? "Сохранить" : "Создать"}
              </Button>
            </Group>
          </Stack>
        </form>
      </Modal>
    </>
  );
}

export default Transactions;
