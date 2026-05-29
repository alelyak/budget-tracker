import { useState, useEffect } from "react";
import {
  Title,
  Text,
  Card,
  Progress,
  Loader,
  Center,
  Alert,
  SimpleGrid,
} from "@mantine/core";
import { IconAlertCircle, IconInfoCircle } from "@tabler/icons-react";
import {
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import api from "../api/axios";

function Dashboard() {
  const [stats, setStats] = useState(null);
  const [expensesByCategory, setExpensesByCategory] = useState([]);
  const [dailyTrend, setDailyTrend] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const currentDate = new Date();
    const currentMonth = `${currentDate.getFullYear()}-${String(currentDate.getMonth() + 1).padStart(2, "0")}`;

    Promise.all([
      api.get(`/dashboard/stats?month=${currentMonth}`),
      api.get(`/charts/expenses-by-category?month=${currentMonth}`),
      api.get(`/charts/daily-trend?month=${currentMonth}`),
    ])
      .then(([statsRes, expensesRes, trendRes]) => {
        setStats(statsRes.data);
        setExpensesByCategory(expensesRes.data || []);
        setDailyTrend(trendRes.data || []);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching dashboard data:", err);
        setError("Не удалось загрузить данные. Убедитесь, что сервер запущен.");
        setLoading(false);
      });
  }, []);

  const COLORS = [
    "#40c057",
    "#f59f00",
    "#fa5252",
    "#228be6",
    "#be4bdb",
    "#0effb7",
    "#5e005a",
    "#e2fd14",
    "#b8ff8b",
    "#bb00ff",
    "#fd1498",
    "#14e6fd",
    "#6614fd",
    "#ff0000",
    "#1d6685",

  ];

  if (loading) {
    return (
      <Center style={{ height: "60vh", flexDirection: "column" }}>
        <Loader size="xl" mb="md" />
        <Text c="dimmed">Загрузка данных...</Text>
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

  const hasData = stats && (stats.total_income > 0 || stats.total_expense > 0);

  return (
    <>
      <Title order={2} mb="md">
        📊 Дашборд
      </Title>

      <SimpleGrid cols={3} spacing="md" mb="xl">
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Text size="sm" c="dimmed">
            Баланс
          </Text>
          <Text size="28px" fw={700} c="blue">
            {(stats?.balance || 0).toLocaleString()} ₽
          </Text>
        </Card>
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Text size="sm" c="dimmed">
            Доходы (месяц)
          </Text>
          <Text size="28px" fw={700} c="green">
            +{(stats?.total_income || 0).toLocaleString()} ₽
          </Text>
        </Card>
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Text size="sm" c="dimmed">
            Расходы (месяц)
          </Text>
          <Text size="28px" fw={700} c="red">
            -{(stats?.total_expense || 0).toLocaleString()} ₽
          </Text>
        </Card>
      </SimpleGrid>

      <SimpleGrid cols={2} spacing="md" mb="xl">
        {/* Круговая диаграмма с легендой */}
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Title order={4} mb="md" ta="center">
            Расходы по категориям
          </Title>
          {expensesByCategory.length === 0 ? (
            <Center style={{ height: 250 }}>
              <Text c="dimmed" ta="center">
                Нет данных о расходах
              </Text>
            </Center>
          ) : (
            <>
              <ResponsiveContainer width="100%" height={250}>
                <PieChart>
                  <Pie
                    data={expensesByCategory}
                    cx="50%"
                    cy="50%"
                    labelLine={false}
                    outerRadius={80}
                    dataKey="amount"
                    nameKey="category_name"
                  >
                    {expensesByCategory.map((entry, index) => (
                      <Cell
                        key={`cell-${index}`}
                        fill={COLORS[index % COLORS.length]}
                      />
                    ))}
                  </Pie>
                  <Tooltip
                    formatter={(value) => `${value.toLocaleString()} ₽`}
                  />
                </PieChart>
              </ResponsiveContainer>

              {/* Легенда под диаграммой */}
              <div
                style={{
                  marginTop: "16px",
                  display: "flex",
                  flexWrap: "wrap",
                  gap: "12px",
                  justifyContent: "center",
                }}
              >
                {expensesByCategory.map((entry, index) => (
                  <div
                    key={entry.category_name}
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: "6px",
                    }}
                  >
                    <div
                      style={{
                        width: 12,
                        height: 12,
                        borderRadius: "50%",
                        backgroundColor: COLORS[index % COLORS.length],
                      }}
                    />
                    <Text size="xs">
                      {entry.category_name}: {entry.percentage.toFixed(0)}% (
                      {entry.amount.toLocaleString()} ₽)
                    </Text>
                  </div>
                ))}
              </div>
            </>
          )}
        </Card>

        {/* Столбчатая диаграмма */}
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Title order={4} mb="md" ta="center">
            Динамика за месяц
          </Title>
          {dailyTrend.length === 0 ? (
            <Center style={{ height: 250 }}>
              <Text c="dimmed" ta="center">
                Нет данных за месяц
              </Text>
            </Center>
          ) : (
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={dailyTrend}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" tick={{ fontSize: 12 }} />
                <YAxis tickFormatter={(value) => `${value / 1000}K`} />
                <Tooltip formatter={(value) => `${value.toLocaleString()} ₽`} />
                <Legend />
                <Bar dataKey="income" name="Доходы" fill="#40c057" />
                <Bar dataKey="expense" name="Расходы" fill="#fa5252" />
              </BarChart>
            </ResponsiveContainer>
          )}
        </Card>
      </SimpleGrid>

      {/* Прогресс-бары бюджета */}
      {stats?.category_budgets?.length > 0 && (
        <Card shadow="sm" padding="lg" radius="md" withBorder mb="xl">
          <Title order={3} mb="md">
            📊 Бюджет по категориям
          </Title>
          {stats.category_budgets.map((cat) => (
            <div key={cat.id} style={{ marginBottom: "20px" }}>
              <div
                style={{
                  display: "flex",
                  justifyContent: "space-between",
                  marginBottom: "8px",
                }}
              >
                <Text fw={500}>{cat.name}</Text>
                <Text size="sm" c={cat.is_exceeded ? "red" : "dimmed"}>
                  {cat.spent.toLocaleString()} /{" "}
                  {cat.monthly_limit.toLocaleString()} ₽
                </Text>
              </div>
              <Progress
                value={Math.min(cat.percent, 100)}
                color={
                  cat.is_exceeded
                    ? "red"
                    : cat.percent > 80
                      ? "yellow"
                      : "green"
                }
                size="lg"
                radius="md"
              />
              {cat.is_exceeded && (
                <Text size="xs" c="red" mt={4}>
                  ⚠️ Превышение на{" "}
                  {(cat.spent - cat.monthly_limit).toLocaleString()} ₽
                </Text>
              )}
            </div>
          ))}
        </Card>
      )}

      {/* Последние транзакции */}
      {stats?.recent_transactions?.length > 0 && (
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Title order={3} mb="md">
            🕒 Последние операции
          </Title>
          <div>
            {stats.recent_transactions.slice(0, 5).map((tx) => (
              <div
                key={tx.id}
                style={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  padding: "12px 0",
                  borderBottom: "1px solid #e9ecef",
                }}
              >
                <div>
                  <Text size="sm" c="dimmed">
                    {tx.date?.split("T")[0]}
                  </Text>
                  <Text fw={500}>{tx.category_name}</Text>
                  <Text size="xs" c="dimmed">
                    {tx.note || "—"}
                  </Text>
                </div>
                <Text
                  size="md"
                  fw={600}
                  c={tx.category_type === "income" ? "green" : "red"}
                >
                  {tx.category_type === "income" ? "+" : "-"}
                  {tx.amount.toLocaleString()} ₽
                </Text>
              </div>
            ))}
          </div>
        </Card>
      )}

      {!hasData && (
        <Card shadow="sm" padding="xl" radius="md" withBorder>
          <Center style={{ flexDirection: "column", padding: "40px 20px" }}>
            <IconInfoCircle size={48} color="#868e96" stroke={1.5} />
            <Title order={3} mt="md" c="dimmed">
              Нет данных
            </Title>
            <Text c="dimmed" ta="center" mt="sm">
              Добавьте категории и транзакции, чтобы увидеть статистику и
              графики
            </Text>
          </Center>
        </Card>
      )}
    </>
  );
}

export default Dashboard;
