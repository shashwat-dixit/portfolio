'use client'
import { motion } from 'framer-motion'

/*
  Add carousel to showcase projects
  And add a skills tab to showcase all the skills with their levels.
*/

export default function WorkPage() {
  return (
    <motion.div
      initial={{ y: -50, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.6, ease: "easeOut" }}
      className="flex justify-center text-3xl text-nowrap absolute left-1/2 right-1/2 top-28 transform -translate-x-1/2"
    >
      {/* Card elements to showcase work */}
      Work Page
    </motion.div>
  )
}
