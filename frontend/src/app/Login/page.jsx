"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "../context/Authcontext";

export default function Login() {
  const { setAuth } = useAuth(); // ✅ FIXED: hook used INSIDE component
  const router = useRouter();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage("");
    setLoading(true);

    try {
      const res = await fetch("http://localhost:8080/api/v1/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
        credentials: "include",
      });

      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.error || "Login failed");
      }

      const data = await res.json();
      console.log("Login Response:", data);

      // Save tokens
      localStorage.setItem("access_token", data.access_token);
      localStorage.setItem("user_id", data.user_id);
      localStorage.setItem("user_email", email);

      // Update Auth Context globally
      setAuth({
        accessToken: data.access_token,
        userId: data.user_id,
        loggedIn: true,
      });

      setMessage("Login successful! Redirecting...");
      router.push("/");

    } catch (err) {
      setMessage(err.message || "Something went wrong");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-screen h-screen flex justify-center items-center bg-transparent bg-gradient-to-br from-green-500 to-green-800">
      <div className="bg-green-800/90 border-2 border-green-700 rounded-xl shadow-xl px-8 py-6 max-w-md w-full">
        <div className="text-center mb-6">
          <p className="text-2xl font-semibold text-white">Login</p>
          <p className="text-sm text-green-100 mt-1">Access your dashboard</p>
        </div>

        <form className="space-y-4 text-left" onSubmit={handleSubmit}>
          <div className="flex flex-col gap-2">
            <label className="text-sm text-green-100" htmlFor="email">
              Enter Email
            </label>
            <input
              id="email"
              type="email"
              placeholder="Enter your email"
              className="border border-green-300/60 rounded-md bg-white/90 px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-green-500"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          <div className="flex flex-col gap-2">
            <label className="text-sm text-green-100" htmlFor="password">
              Enter Password
            </label>
            <input
              id="password"
              type="password"
              placeholder="Enter your password"
              className="border border-green-300/60 rounded-md bg-white/90 px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-green-500"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          {message && (
            <div className="text-sm text-green-100 bg-green-900/60 border border-green-700 rounded px-3 py-2">
              {message}
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full mt-2 bg-green-600 hover:bg-green-500 disabled:bg-green-700 text-white font-semibold py-2 rounded-md transition-colors duration-150"
          >
            {loading ? "Signing in..." : "Login"}
          </button>
        </form>
      </div>
    </div>
  );
}
