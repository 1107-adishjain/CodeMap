"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export default function Signup() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [timeout , settimeout] = useState("");
  const router = useRouter();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage("");
    if (password !== confirm) {
      setMessage("Passwords do not match");
      return;
    }
    setLoading(true);
    try {
      const res = await fetch("http://localhost:8080/api/v1/signup", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(data.error || "Signup failed");
      setMessage("Signup successful! You can now log in.");
      setEmail("");
      setPassword("");
      setConfirm("");
if(res.status == 201){
    router.push("/Login")
}
      
    } catch (err) {
      setMessage(err.message || "Something went wrong");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-screen h-screen flex justify-center items-center  bg-gradient-to-br from-green-500 to-green-800">
      <div className="bg-green-800/90 border-2 border-green-700 rounded-xl shadow-xl px-8 py-6 max-w-md w-full">
        <div className="text-center mb-6">
          <p className="text-2xl font-semibold text-white">Sign Up</p>
          <p className="text-sm text-green-100 mt-1">Create your account</p>
        </div>

        <form className="space-y-4 text-left" onSubmit={handleSubmit}>
          <div className="flex flex-col gap-2">
            <label className="text-sm text-green-100" htmlFor="email">
              Enter Email
            </label>
            <input
              id="email"
              type="email"
              placeholder="user@gmail.com"
              className="border border-green-300/60 rounded-md bg-white/90 px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-green-500"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          <div className="flex flex-col gap-2">
            <label className="text-sm text-green-100" htmlFor="password">
              Enter Password (min 8 chars)
            </label>
            <input
              id="password"
              type="password"
              placeholder="••••••••"
              className="border border-green-300/60 rounded-md bg-white/90 px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-green-500"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          <div className="flex flex-col gap-2">
            <label className="text-sm text-green-100" htmlFor="confirm">
              Confirm Password
            </label>
            <input
              id="confirm"
              type="password"
              placeholder="••••••••"
              className="border border-green-300/60 rounded-md bg-white/90 px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-green-500"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              required
            />
          </div>

          {message && (
            <div
              className={`text-sm px-3 py-2 rounded border ${
                message.toLowerCase().includes("success")
                  ? "bg-green-900/60 border-green-700 text-green-100"
                  : "bg-red-900/60 border-red-700 text-red-100"
              }`}
            >
              {message}
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full mt-2 bg-green-600 hover:bg-green-500 disabled:bg-green-700 text-white font-semibold py-2 rounded-md transition-colors duration-150"
          >
            {loading ? "Signing up..." : "Sign Up"}
          </button>
        </form>
      </div>
    </div>
  );
}