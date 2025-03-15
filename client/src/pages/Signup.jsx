import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import axios from "axios";

function Signup() {
  const navigate = useNavigate();

  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const handleUserSign = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");
    
    let fullName = firstName + " " + lastName;
    let data = {
      Name: fullName,
      Email: email,
      Password: password
    };
    
    try {
      let response = await axios.post("http://localhost:8080/signup", data);
      console.log("success");
      setEmail("");
      setFirstName("");
      setLastName("");
      setPassword("");
      navigate("/Home");
    } catch (e) {
      console.log("error : ", e);
      setError("Failed to sign up. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex justify-center items-center bg-gradient-to-br from-blue-100 to-purple-100 px-4 sm:px-6 lg:px-8 py-12 min-h-screen">
      <div className="flex bg-white shadow-2xl rounded-xl w-full max-w-4xl overflow-hidden">
        {/* Left side - Image */}
        <div className="hidden md:block relative md:w-1/2">
          <div className="z-10 absolute inset-0 bg-gradient-to-br from-blue-600/80 to-purple-800/90"></div>
          <img
            src="https://images.unsplash.com/photo-1634363657957-d91ac22d230a?fm=jpg&q=60&w=3000&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA=="
            alt="Signup Illustration"
            className="absolute inset-0 w-full h-full object-cover"
          />
          <div className="z-20 absolute inset-0 flex flex-col justify-center items-center p-8">
            <h1 className="mb-6 font-bold text-white text-3xl">Welcome!</h1>
            <p className="mb-8 text-white text-lg text-center">Join our community and discover amazing features.</p>
            <div className="bg-white/20 backdrop-blur-sm p-4 rounded-lg">
              <p className="text-white text-sm italic">"The best investment you can make is in yourself."</p>
            </div>
          </div>
        </div>

        {/* Right side - Form */}
        <div className="p-8 w-full md:w-1/2">
          <div className="mb-6 text-center">
            <h2 className="font-extrabold text-gray-900 text-3xl">Create Your Account</h2>
            <p className="mt-2 text-gray-600 text-sm">
              Start your journey with us today
            </p>
          </div>

          {error && (
            <div className="bg-red-100 mb-4 p-3 border border-red-400 rounded text-red-700">
              {error}
            </div>
          )}

          <form onSubmit={handleUserSign} className="space-y-5">
            <div className="gap-4 grid grid-cols-1 sm:grid-cols-2">
              <div>
                <label htmlFor="firstName" className="block mb-1 font-medium text-gray-700 text-sm">
                  First Name
                </label>
                <input
                  id="firstName"
                  type="text"
                  required
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  className="px-4 py-3 border border-gray-300 focus:border-blue-500 rounded-lg focus:ring-2 focus:ring-blue-500 w-full transition-all"
                  placeholder="John"
                />
              </div>
              <div>
                <label htmlFor="lastName" className="block mb-1 font-medium text-gray-700 text-sm">
                  Last Name
                </label>
                <input
                  id="lastName"
                  type="text"
                  required
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  className="px-4 py-3 border border-gray-300 focus:border-blue-500 rounded-lg focus:ring-2 focus:ring-blue-500 w-full transition-all"
                  placeholder="Doe"
                />
              </div>
            </div>

            <div>
              <label htmlFor="email" className="block mb-1 font-medium text-gray-700 text-sm">
                Email Address
              </label>
              <input
                id="email"
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="px-4 py-3 border border-gray-300 focus:border-blue-500 rounded-lg focus:ring-2 focus:ring-blue-500 w-full transition-all"
                placeholder="john@example.com"
              />
            </div>

            <div>
              <label htmlFor="password" className="block mb-1 font-medium text-gray-700 text-sm">
                Password
              </label>
              <input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="px-4 py-3 border border-gray-300 focus:border-blue-500 rounded-lg focus:ring-2 focus:ring-blue-500 w-full transition-all"
                placeholder="••••••••"
              />
              <p className="mt-1 text-gray-500 text-xs">
                Password should be at least 8 characters
              </p>
            </div>

            <div className="flex items-center">
              <input
                id="newsletter"
                type="checkbox"
                className="border-gray-300 rounded focus:ring-blue-500 w-4 h-4 text-blue-600"
              />
              <label htmlFor="newsletter" className="block ml-2 text-gray-700 text-sm">
                Subscribe to our newsletter
              </label>
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className="flex justify-center bg-gradient-to-r from-blue-600 hover:from-blue-700 to-blue-700 hover:to-blue-800 shadow-sm px-4 py-3 border border-transparent rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 w-full font-medium text-white text-lg transition-all"
            >
              {isLoading ? "Signing up..." : "Sign Up"}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-gray-600 text-sm">
              Already have an account?{" "}
              <Link
                to="/Login"
                className="font-medium text-blue-600 hover:text-blue-500 transition-colors"
              >
                Sign in instead
              </Link>
            </p>
          </div>

          <div className="mt-6">
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="border-gray-300 border-t w-full"></div>
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="bg-white px-2 text-gray-500">Or continue with</span>
              </div>
            </div>

            <div className="gap-3 grid grid-cols-2 mt-6">
              <button className="inline-flex justify-center bg-white hover:bg-gray-50 shadow-sm px-4 py-2 border border-gray-300 rounded-lg w-full font-medium text-gray-700 text-sm">
                <svg className="w-5 h-5" fill="#4285F4" viewBox="0 0 24 24">
                  <path d="M12.545,10.239v3.821h5.445c-0.712,2.315-2.647,3.972-5.445,3.972c-3.332,0-6.033-2.701-6.033-6.032s2.701-6.032,6.033-6.032c1.498,0,2.866,0.549,3.921,1.453l2.814-2.814C17.503,2.988,15.139,2,12.545,2C7.021,2,2.543,6.477,2.543,12s4.478,10,10.002,10c8.396,0,10.249-7.85,9.426-11.748L12.545,10.239z"/>
                </svg>
              </button>
              <button className="inline-flex justify-center bg-white hover:bg-gray-50 shadow-sm px-4 py-2 border border-gray-300 rounded-lg w-full font-medium text-gray-700 text-sm">
                <svg className="w-5 h-5" fill="#3b5998" viewBox="0 0 24 24">
                  <path d="M22,12c0-5.52-4.48-10-10-10S2,6.48,2,12c0,4.84,3.44,8.87,8,9.8V15H8v-3h2V9.5C10,7.57,11.57,6,13.5,6H16v3h-2 c-0.55,0-1,0.45-1,1v2h3v3h-3v6.95C18.05,21.45,22,17.19,22,12z"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Signup;