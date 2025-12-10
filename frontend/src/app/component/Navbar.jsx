"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";

const navItems = [
  { label: "Home", path: "/" },
  { label: "About", path: "/about" },
  { label: "Services", path: "/services" },
  { label: "Contact", path: "/contact" },
];

const Navbar = () => {
  const router = useRouter();
  const [isOpen, setIsOpen] = useState(false);
  const [user, setUser] = useState(null);

  useEffect(() => {
    // Check if user is logged in
    const email = localStorage.getItem("user_email");
    // const token = localStorage.getItem("access_token");
    if (email) {
      setUser({ email });
    }
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("user_email");
    // localStorage.removeItem("access_token");
    // localStorage.removeItem("user_id");
    setUser(null);
    router.push("/");
  };

  const toggleMenu = () => {
    setIsOpen(!isOpen);
  };

  return (
    <nav className="fixed top-0 left-0 w-full bg-gradient-to-r from-green-700 to-teal-600 text-white shadow-md z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <div className="flex-shrink-0">
            <h1 className="text-2xl font-bold cursor-pointer" onClick={() => router.push("/")}>
              CodeMap
            </h1>
          </div>

          {/* Desktop Menu */}
          <div className="hidden md:flex space-x-6 items-center">
            {navItems.map((item) => (
              <button
                key={item.label}
                onClick={() => router.push(item.path)}
                className="hover:text-teal-300 transition-colors duration-300"
              >
                {item.label}
              </button>
            ))}

            {/* Auth Buttons or User Email */}
            {user ? (
              <div className="flex items-center space-x-4">
                <span className="text-sm bg-green-800 px-3 py-1 rounded-full">
                  {user.email}
                </span>
                <button
                  onClick={handleLogout}
                  className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md transition-colors duration-300"
                >
                  Logout
                </button>
              </div>
            ) : (
              <>
                <button
                  onClick={() => router.push("/Login")}
                  className="bg-white text-green-700 hover:bg-gray-100 px-4 py-2 rounded-md transition-colors duration-300"
                >
                  Login
                </button>
                <button
                  onClick={() => router.push("/Signup")}
                  className="bg-teal-500 hover:bg-teal-600 text-white px-4 py-2 rounded-md transition-colors duration-300"
                >
                  Register
                </button>
              </>
            )}
          </div>

          {/* Mobile Menu Button */}
          <div className="md:hidden">
            <button
              onClick={toggleMenu}
              className="text-white focus:outline-none focus:ring-2 focus:ring-teal-300"
            >
              <svg
                className="h-6 w-6"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                {isOpen ? (
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                ) : (
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 6h16M4 12h16m-7 6h7"
                  />
                )}
              </svg>
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Menu */}
      {isOpen && (
        <div className="md:hidden bg-green-800">
          <div className="px-2 pt-2 pb-3 space-y-1">
            {navItems.map((item) => (
              <button
                key={item.label}
                onClick={() => {
                  router.push(item.path);
                  setIsOpen(false);
                }}
                className="block w-full text-left px-3 py-2 rounded-md text-base font-medium hover:bg-green-700 transition-colors duration-300"
              >
                {item.label}
              </button>
            ))}

            {/* Mobile Auth Buttons or User Email */}
            {user ? (
              <div className="space-y-2">
                <div className="px-3 py-2 text-sm bg-green-900 rounded-md">
                  {user.email}
                </div>
                <button
                  onClick={() => {
                    handleLogout();
                    setIsOpen(false);
                  }}
                  className="block w-full text-left px-3 py-2 bg-red-600 hover:bg-red-700 rounded-md transition-colors duration-300"
                >
                  Logout
                </button>
              </div>
            ) : (
              <>
                <button
                  onClick={() => {
                    router.push("/Login");
                    setIsOpen(false);
                  }}
                  className="block w-full text-left px-3 py-2 bg-white text-green-700 hover:bg-gray-100 rounded-md transition-colors duration-300"
                >
                  Login
                </button>
                <button
                  onClick={() => {
                    router.push("/signup");
                    setIsOpen(false);
                  }}
                  className="block w-full text-left px-3 py-2 bg-teal-500 hover:bg-teal-600 rounded-md transition-colors duration-300"
                >
                  Register
                </button>
              </>
            )}
          </div>
        </div>
      )}
    </nav>
  );
};

export default Navbar;