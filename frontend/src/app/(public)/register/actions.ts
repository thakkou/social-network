"use server";

import { signIn } from "~/server/auth";

type RegisterResult = { error: string } | { success: true };

export async function registerUser(
  formData: FormData,
): Promise<RegisterResult> {
  const firstname = formData.get("firstname")?.toString().trim();
  const lastname = formData.get("lastname")?.toString().trim();
  const email = formData.get("email")?.toString().trim();
  const password = formData.get("password")?.toString();
  const birthDate = formData.get("birthDate")?.toString();

  // Optional fields
  const nickname = formData.get("nickname")?.toString().trim();
  const aboutme = formData.get("aboutme")?.toString().trim();
  const avatar = formData.get("avatar");

  if (!firstname || !lastname || !email || !password || !birthDate) {
    return { error: "All required fields must be provided." };
  }

  const goFormData = new FormData();

  // Required fields
  goFormData.append("firstname", firstname);
  goFormData.append("lastname", lastname);
  goFormData.append("email", email);
  goFormData.append("password", password);
  goFormData.append("birthDate", birthDate);

  // Optional text fields
  if (nickname) {
    goFormData.append("nickname", nickname);
  }

  if (aboutme) {
    goFormData.append("aboutme", aboutme);
  }

  // Optional avatar
  if (
    avatar instanceof File &&
    avatar.size > 0 &&
    avatar.name !== "undefined"
  ) {
    goFormData.append("avatar", avatar, avatar.name);
  }

  // Debug: verify what is actually sent
  console.log("Sending FormData:");
  for (const [key, value] of goFormData.entries()) {
    console.log(key, value);
  }

  const res = await fetch(
    `${process.env.GO_BACKEND_URL}/api/register`,
    {
      method: "POST",
      body: goFormData,
    },
  );

  if (!res.ok) {
    const data = await res.json().catch(() => null);
    console.error(data);
    return { error: data?.error ?? "Registration failed" };
  }

  // Uncomment if you want to auto-login after registration
  /*
  await signIn("credentials", {
    email,
    password,
    redirect: false,
  });
  */

  return { success: true };
}