// Mock request data — mirrors the API shape from
// dynasty/server/handlers/requests/transport/dto.go
window.GUARD_DATA = {
  reqTypes: {
    1: { key: "guest",    ua: "Гість",              color: "var(--rtype-guest)" },
    2: { key: "taxi",     ua: "Таксі",              color: "var(--rtype-taxi)" },
    3: { key: "delivery", ua: "Доставка",           color: "var(--rtype-delivery)" },
    4: { key: "cargo",    ua: "37-Б Розвантаження", color: "var(--rtype-cargo)" },
  },
  requests: [
    {
      id: 4821, rtype: 3, status: "new", time: 1715870520,
      address: "37-В", apartment: 142,
      user_name: "Олена К.", phone: "+380 67 123 45 67",
      description: "Кур'єр Нової пошти, замовлення на ім'я Олена. Велика коробка — залиште біля консьєржа якщо мене немає.",
      images: [{ thumb: "", img: "" }],
    },
    {
      id: 4820, rtype: 1, status: "new", time: 1715869200,
      address: "37-В", apartment: 142,
      user_name: "Олена К.", phone: "+380 67 123 45 67",
      description: "О 19:30 прийде Тарас, мій брат. Пропустіть, будь ласка.",
      images: [],
    },
    {
      id: 4819, rtype: 2, status: "new", time: 1715864100,
      address: "37-Б", apartment: 58,
      user_name: "Богдан М.", phone: "+380 50 987 65 43",
      description: "Виклик таксі на 17:00, Bolt, біла Toyota AA 4521 KX.",
      images: [],
    },
    {
      id: 4818, rtype: 4, status: "new", time: 1715855700,
      address: "37-Б", apartment: 12,
      user_name: "Андрій Л.", phone: "+380 63 444 33 22",
      description: "Розвантаження меблів IKEA, фургон ~20 хв.",
      images: [{ thumb: "", img: "" }, { thumb: "", img: "" }],
    },
    {
      id: 4817, rtype: 3, status: "closed", time: 1715842800,
      address: "37-В", apartment: 207,
      user_name: "Марина Г.", phone: "+380 96 111 22 33",
      description: "Glovo замовлення №AB-2241.",
      images: [],
    },
    {
      id: 4816, rtype: 1, status: "closed", time: 1715839200,
      address: "37-В", apartment: 304,
      user_name: "Сергій В.", phone: "+380 67 555 11 22",
      description: "Майстер з кондиціонерів, прізвище Шевчук.",
      images: [],
    },
  ],
};
