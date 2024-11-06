'use client'
import ProjectShowcase from '@/components/ProjectShowcase'
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
      className=""
    >
      <ProjectShowcase />
    </motion.div>
  )
}
