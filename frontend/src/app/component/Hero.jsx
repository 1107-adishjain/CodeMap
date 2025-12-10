import { Infinity } from "lucide-react";
import Link from "next/link";

export default function Hero() {
  return (
    <>
      <div className=" grid grid-cols-2 gap-4 w-full h-screen px-20">
        <div className=" mt-20">
          <h1 className="lg:text-6xl md:text-4xl sm:text-2xl font-bold mb-8">Welcome to <span className="text-teal-400">CODEmap</span></h1>
          <p className="text-2xl text-gray-300">
            Your ultimate platform for visualizing your code architecture.
          </p>
          <p className="text-lg text-gray-300 mt-10">
            Transform your codebase into interactive visual maps. Analyze relationships, understand dependencies,<br /> and explore your software architecture like never before
          </p>


          <div><Link href="/dashboard">
            <button className="mt-6 px-10 py-5 bg-teal-500 text-2xl hover:bg-teal-600 text-white rounded-lg shadow-md transition duration-300 mr-4">
              Start Analyze Code</button>
          </Link>

            <button className="mt-6 px-10 py-5 bg-gray-700 text-2xl hover:bg-gray-600 text-white rounded-lg shadow-md transition duration-300">Learn More</button>
          </div>

          <div className="mt-10 grid lg:grid-cols-3 md:cols-3 sm:grid-cols-1 gap-4">
            <div className=" p-4 rounded-lg shadow-md transition-all delay-100 duration-200 hover:scale-105 hover:shadow-teal-700 text-center">
              <h2 className="text-3xl font-bold mb-2">15+</h2>
              <p className="text-gray-300">Programming languages supported</p>
            </div>
            <div className=" p-4 rounded-lg shadow-md text-center transition-all delay-100 duration-200 hover:scale-105 hover:shadow-teal-700 ">
              <h2 className="text-3xl font-bold mb-2 text-center ">AST</h2>
              <p className="text-gray-300 text-center">Abstract Syntax Tree supported.</p>
            </div>
            <div className="p-4 rounded-lg shadow-md text-center transition-all delay-100 duration-200 hover:scale-105 hover:shadow-teal-700">
              <h2 className="text-3xl font-bold mb-2 flex justify-center"><Infinity size={36} /></h2>
              <p className="text-gray-300">Code Relationships.</p>
            </div>
          </div>
        </div>
        <div>
          <h1>second part </h1>
        </div>
      </div>

    </>

  );
}