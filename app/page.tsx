"use client"
import SignIn from "@/components/sign-in";
import { motion } from "framer-motion";


export default function Home() {
  return (
    <div className="">
      <motion.h1
        initial={{ y: -50, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
        className="flex justify-center text-3xl text-nowrap absolute left-1/2 right-1/2 top-28 transform -translate-x-1/2"
      >
        Hi, I am Shashwat Dixit
      </motion.h1>
      <SignIn />
    </div>
  );
}