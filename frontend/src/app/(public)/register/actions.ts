"use server";

import { signIn } from "~/server/auth";

type RegisterResult = { error: string } | { success: true };

export async function registerUser(formData: FormData): Promise<RegisterResult> {
  const firstname = formData.get("firstname") as string;
  const lastname = formData.get("lastname") as string;
  const email = formData.get("email") as string;
  const password = formData.get("password") as string;
  const birthDate = formData.get("birthDate") as string;
  const nickname = formData.get("nickname") as string;
  const aboutme = formData.get("aboutme") as string;
  const avatar = formData.get("avatar") as File; //  | null

  if (!firstname || !lastname || !email || !password || !birthDate) {
    return { error: "All fields are required" };
  }

  const goFormData = new FormData();
  goFormData.append("firstname", firstname);
  goFormData.append("lastname", lastname);
  goFormData.append("email", email);
  goFormData.append("password", password);
  goFormData.append("birthDate", birthDate);
  goFormData.append("nickname", nickname); // opt
  goFormData.append("aboutme", aboutme); // opt
  // if (avatar && avatar.size > 0) {
    goFormData.append("avatar", avatar, avatar.name); // File objects work directly with fetch's FormData
  // }

  const res = await fetch(`${process.env.GO_BACKEND_URL}/api/register`, {
    method: "POST",
    // headers: { "Content-Type": "application/json" },
    body: goFormData, // JSON.stringify({ firstname, lastname, email, password, birthDate, nickname, aboutme }), // avatar
    // NOTE: do NOT set Content-Type manually — fetch sets the correct
    // multipart boundary header automatically when body is FormData
  });

  if (!res.ok) {
    const data = await res.json().catch(() => null);
    console.log(data) // maybe .data field
    return { error: data?.error ?? "Registration failed" };
  }

  // auto-login right after successful registration, reusing your
  // existing Credentials provider (same authorize() flow as normal login)
//   await signIn("credentials", {
//     email,
//     password,
//     redirect: false,
//   });
    // send to login instead

  return { success: true };
}