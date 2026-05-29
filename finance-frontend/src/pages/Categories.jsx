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
  Alert,
  Badge,
  Box,
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

function Categories() {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [modalOpened, setModalOpened] = useState(false);
  const [editingId, setEditingId] = useState(null);
  // Состояния для модального окна подтверждения удаления
  const [deleteModalOpened, setDeleteModalOpened] = useState(false);
  const [deletingCategory, setDeletingCategory] = useState(null);

  const form = useForm({
    initialValues: {
      name: "",
      type: "expense",
      monthly_limit: 0,
    },
    validate: {
      name: (val) => (!val ? "Введите название категории" : null),
      monthly_limit: (val) =>
        val < 0 ? "Лимит не может быть отрицательным" : null,
    },
  });

  const fetchCategories = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await api.get("/categories");
      setCategories(res.data || []);
    } catch (error) {
      setError(
        "Не удалось загрузить категории. Проверьте подключение к серверу.",
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCategories();
  }, []);

  const handleSubmit = async (values) => {
    try {
      const data = {
        name: values.name,
        type: values.type,
        monthly_limit: parseFloat(values.monthly_limit) || 0,
      };

      if (editingId) {
        await api.put(`/categories/${editingId}`, data);
        notifications.show({
          title: "Успешно",
          message: "Категория обновлена",
          color: "green",
          icon: <IconCheck size={16} />,
          autoClose: 3000,
        });
      } else {
        await api.post("/categories", data);
        notifications.show({
          title: "Успешно",
          message: "Категория создана",
          color: "green",
          icon: <IconCheck size={16} />,
          autoClose: 3000,
        });
      }
      closeModal();
      fetchCategories();
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

  // Открыть модальное окно подтверждения удаления
  const openDeleteModal = (category) => {
    setDeletingCategory(category);
    setDeleteModalOpened(true);
  };

  // Выполнить удаление
  const confirmDelete = async () => {
    if (!deletingCategory) return;

    try {
      await api.delete(`/categories/${deletingCategory.id}`);
      notifications.show({
        title: "Успешно",
        message: `Категория "${deletingCategory.name}" удалена`,
        color: "green",
        icon: <IconCheck size={16} />,
        autoClose: 3000,
      });
      fetchCategories();
    } catch (error) {
      notifications.show({
        title: "Ошибка",
        message: "Не удалось удалить категорию",
        color: "red",
        icon: <IconX size={16} />,
        autoClose: 4000,
      });
    } finally {
      setDeleteModalOpened(false);
      setDeletingCategory(null);
    }
  };

  const openModal = (category = null) => {
    if (category) {
      setEditingId(category.id);
      form.setValues({
        name: category.name,
        type: category.type,
        monthly_limit: category.monthly_limit,
      });
    } else {
      setEditingId(null);
      form.reset();
      form.setValues({ type: "expense", monthly_limit: 0 });
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
        <Text c="dimmed">Загрузка категорий...</Text>
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

  const expenseCategories = categories.filter((c) => c.type === "expense");
  const incomeCategories = categories.filter((c) => c.type === "income");

  return (
    <>
      <Group justify="space-between" mb="md">
        <Title order={2}>📂 Категории</Title>
        <Button
          leftSection={<IconPlus size={20} />}
          onClick={() => openModal()}
        >
          Добавить категорию
        </Button>
      </Group>

      {/* Категории расходов */}
      <Card shadow="sm" padding="lg" radius="md" withBorder mb="md">
        <Title order={3} mb="md" c="red">
          💸 Расходы
        </Title>

        {expenseCategories.length === 0 ? (
          <Center style={{ flexDirection: "column", padding: "40px 20px" }}>
            <Text c="dimmed">Нет категорий расходов</Text>
            <Button
              leftSection={<IconPlus size={20} />}
              onClick={() => openModal()}
              variant="outline"
              size="xs"
              mt="md"
            >
              Добавить категорию
            </Button>
          </Center>
        ) : (
          <Table striped highlightOnHover>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Название</Table.Th>
                <Table.Th>Месячный лимит</Table.Th>
                <Table.Th style={{ width: 100 }}>Действия</Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {expenseCategories.map((cat) => (
                <Table.Tr key={cat.id}>
                  <Table.Td fw={500}>{cat.name}</Table.Td>
                  <Table.Td>
                    {cat.monthly_limit > 0 ? (
                      <Badge color="blue" size="lg">
                        {cat.monthly_limit.toLocaleString()} ₽
                      </Badge>
                    ) : (
                      <Text c="dimmed" size="sm">
                        Без лимита
                      </Text>
                    )}
                  </Table.Td>
                  <Table.Td>
                    <Group gap="xs">
                      <Button
                        size="xs"
                        variant="subtle"
                        color="blue"
                        onClick={() => openModal(cat)}
                      >
                        <IconEdit size={16} />
                      </Button>
                      <Button
                        size="xs"
                        variant="subtle"
                        color="red"
                        onClick={() => openDeleteModal(cat)}
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

      {/* Категории доходов */}
      <Card shadow="sm" padding="lg" radius="md" withBorder>
        <Title order={3} mb="md" c="green">
          💰 Доходы
        </Title>

        {incomeCategories.length === 0 ? (
          <Center style={{ flexDirection: "column", padding: "40px 20px" }}>
            <Text c="dimmed">Нет категорий доходов</Text>
            <Button
              leftSection={<IconPlus size={20} />}
              onClick={() => openModal()}
              variant="outline"
              size="xs"
              mt="md"
            >
              Добавить категорию
            </Button>
          </Center>
        ) : (
          <Table striped highlightOnHover>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Название</Table.Th>
                <Table.Th style={{ width: 100 }}>Действия</Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {incomeCategories.map((cat) => (
                <Table.Tr key={cat.id}>
                  <Table.Td fw={500}>{cat.name}</Table.Td>
                  <Table.Td>
                    <Group gap="xs">
                      <Button
                        size="xs"
                        variant="subtle"
                        color="blue"
                        onClick={() => openModal(cat)}
                      >
                        <IconEdit size={16} />
                      </Button>
                      <Button
                        size="xs"
                        variant="subtle"
                        color="red"
                        onClick={() => openDeleteModal(cat)}
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

      {/* Модальное окно добавления/редактирования */}
      <Modal
        opened={modalOpened}
        onClose={closeModal}
        title={editingId ? "Редактировать категорию" : "Новая категория"}
        size="md"
      >
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <Stack>
            <TextInput
              label="Название"
              placeholder="Например: Продукты, Транспорт"
              required
              {...form.getInputProps("name")}
            />
            <Select
              label="Тип"
              data={[
                { value: "expense", label: "💸 Расход" },
                { value: "income", label: "💰 Доход" },
              ]}
              required
              {...form.getInputProps("type")}
            />
            {form.values.type === "expense" && (
              <NumberInput
                label="Месячный лимит (₽)"
                placeholder="Оставьте 0 или пустым, если лимита нет"
                min={0}
                step={1000}
                {...form.getInputProps("monthly_limit")}
              />
            )}
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

      {/* Модальное окно подтверждения удаления */}
      <Modal
        opened={deleteModalOpened}
        onClose={() => setDeleteModalOpened(false)}
        title="Подтверждение удаления"
        size="sm"
      >
        <Stack>
          <Box>
            <Text size="md">
              Вы действительно хотите удалить категорию{" "}
              <Text component="span" fw={700} c="red">
                "{deletingCategory?.name}"
              </Text>
              ?
            </Text>
            <Text size="sm" c="dimmed" mt="sm">
              ⚠️ Все транзакции, связанные с этой категорией, также будут
              удалены.
            </Text>
          </Box>
          <Group justify="flex-end" mt="md">
            <Button
              variant="subtle"
              onClick={() => setDeleteModalOpened(false)}
            >
              Отмена
            </Button>
            <Button color="red" onClick={confirmDelete}>
              Удалить
            </Button>
          </Group>
        </Stack>
      </Modal>
    </>
  );
}

export default Categories;
