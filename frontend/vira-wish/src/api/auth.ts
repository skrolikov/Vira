export async function login(username: string, password: string) {
  const res = await fetch("/api/id/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  });

  if (!res.ok) throw new Error("Ошибка авторизации");
  return await res.json();
}
