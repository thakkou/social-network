"use client";

// import { signIn } from "next-auth/react";
import { useState, useTransition } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

import { registerUser } from "./actions";

export default function Register() {
  const [firstname, setFirstname] = useState("");
  const [lastname, setLastname] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  // A birth date is usually stored as a string (YYYY-MM-DD) when using an HTML date input.
  const [birthDate, setBirthDate] = useState("");
  // optional
  const [nickname, setNickname] = useState("");
  const [aboutme, setAboutme] = useState("");
  // + avatar

  const [error, setError] = useState("");
  const [preview, setPreview] = useState<string | null>(null);
  const [isPending, startTransition] = useTransition();
  const router = useRouter();

  const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return setPreview(null);

    if (file.size > 1 * 1024 * 1024) {
      setError("Image must be under 1MB");
      e.target.value = "";
      return;
    }

    setPreview(URL.createObjectURL(file)); // instant local preview, no upload yet
  };

  const handleSubmit = async (formData: FormData) => {
    startTransition(async () => {
      console.log(firstname, lastname, email, password, birthDate)
      const result = await registerUser(formData);
      console.log(result)

      if ("error" in result) {
        setError(result.error);
      } else {
        router.push("/");
        router.refresh(); // make sure server components re-read the new session
      }
    });
  //   const res = await signIn("credentials", {
  //     identifier,
  //     password,
  //     redirect: false,
  //   });

  //   if (res?.error) {
  //     setError("Invalid email/nickname or password");
  //   } else {
  //     router.push("/");
  //   }
  };

  // accept="image/png, image/jpeg, image/webp"

  return (
    <form style={{ padding:'2rem', display: 'flex', flexDirection: 'column', justifyContent: 'center' }} id="register-panel" action={handleSubmit}>
      <p className="sec-label" style={{ padding:0, marginBottom:'1rem' }}>create account</p>
      {error && <p>{error}</p>}
      <div style={{ display: 'grid', gridTemplateColumns:'1fr 1fr', gap:'8px' }}>
        <div className="form-row"><span className="form-label">first name *</span><input className="inp" name="firstname" type="text" value={firstname} onChange={(e) => setFirstname(e.target.value)} placeholder="Amir" /></div>
        <div className="form-row"><span className="form-label">last name *</span><input className="inp" name="lastname" type="text" value={lastname} onChange={(e) => setLastname(e.target.value)} placeholder="Kader" /></div>
      </div>
      <div className="form-row"><span className="form-label">email *</span><input className="inp" name="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="you@example.com" /></div>
      <div className="form-row"><span className="form-label">password *</span><input className="inp" name="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••" /></div>
      <div className="form-row"><span className="form-label">date of birth *</span><input className="inp" name="birthDate" type="date" value={birthDate} onChange={(e) => setBirthDate(e.target.value)}/></div>
      <div className="form-row"><span className="form-label">nickname <span className="tag tag-gray">optional</span></span><input className="inp" name="nickname" type="text" value={nickname} onChange={(e) => setNickname(e.target.value)} placeholder="@handle" /></div>
      <div className="form-row"><span className="form-label">about me <span className="tag tag-gray">optional</span></span><textarea className="inp" name="aboutme" value={aboutme} onChange={(e) => setAboutme(e.target.value)} style={{ height:'50px', resize: 'none' }} placeholder="a few words..."></textarea></div>
      <div className="form-row">
        <span className="form-label">avatar <span className="tag tag-gray">optional</span></span>
        {preview && (
          <img src={preview} alt="Avatar preview" width={80} height={80} style={{ borderRadius: "50%" }} />
        )}
        <input id="avatar" name="avatar" type="file" accept=".jpg,.jpeg,.png,.gif" onChange={handleAvatarChange} hidden />
        <label htmlFor="avatar" style={{ border:'0.5px dashed grey', padding:'10px', textAlign: 'center', fontSize:'11px', color:'#e8e4dc', cursor: 'pointer' }}>
          <i className="ti ti-upload" style={{ fontSize:'16px', display: 'block', marginBottom:'4px' }} aria-hidden="true"></i>
          jpeg / png / gif
        </label>
      </div>
      <button className="btn btn-p" style={{ width:'100%', marginTop:'6px' }} type="submit" disabled={isPending}>
        {isPending ? "Creating account..." : "Create account →"}
      </button>
      <p style={{ fontSize:'11px', color:'var(--color-text-tertiary)', marginTop:'10px', textAlign: 'center' }}>already have an account? <Link style={{ color:'#D4537E', cursor: 'pointer' }} href="/login">login</Link></p>
    </form>
  );
}