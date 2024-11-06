'use client'

import { useState } from 'react'
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Github, ExternalLink } from 'lucide-react'

interface Project {
    id: number
    title: string
    description: string
    image: string
    githubLink: string
    deploymentLink: string
    techStack: string[]
}

const projects: Project[] = [
    {
        id: 1,
        title: "Project 1",
        description: "A brief description of Project 1. This project showcases my skills in React and TypeScript.",
        image: "/placeholder.svg?height=200&width=300",
        githubLink: "https://github.com/yourusername/project1",
        deploymentLink: "https://project1.example.com",
        techStack: ["React", "TypeScript", "Tailwind CSS"]
    },
    {
        id: 2,
        title: "Project 2",
        description: "Project 2 is a full-stack application built with Next.js and MongoDB.",
        image: "/placeholder.svg?height=200&width=300",
        githubLink: "https://github.com/yourusername/project2",
        deploymentLink: "https://project2.example.com",
        techStack: ["Next.js", "MongoDB", "Express"]
    },
    {
        id: 3,
        title: "Project 3",
        description: "An innovative mobile app developed using React Native and Firebase.",
        image: "/placeholder.svg?height=200&width=300",
        githubLink: "https://github.com/yourusername/project3",
        deploymentLink: "https://project3.example.com",
        techStack: ["React Native", "Firebase", "Expo"]
    },
    {
        id: 4,
        title: "Project 4",
        description: "A machine learning project that uses Python and TensorFlow for image recognition.",
        image: "/placeholder.svg?height=200&width=300",
        githubLink: "https://github.com/yourusername/project4",
        deploymentLink: "https://project4.example.com",
        techStack: ["Python", "TensorFlow", "Scikit-learn"]
    }
]

export default function ProjectShowcase() {
    const [hoveredProject, setHoveredProject] = useState<number | null>(null)

    return (
        <div className="container mx-auto px-4 py-8 mt-20">
            <h1 className="text-3xl font-bold mb-6 text-center">My Projects</h1>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {projects.map((project) => (
                    <Card
                        key={project.id}
                        className="overflow-hidden transition-all duration-300 transform hover:scale-105"
                        onMouseEnter={() => setHoveredProject(project.id)}
                        onMouseLeave={() => setHoveredProject(null)}
                    >
                        <div className="relative h-48">
                            <img
                                src={project.image}
                                alt={project.title}
                                className="w-full h-full object-cover"
                            />
                            {hoveredProject === project.id && (
                                <div className="absolute inset-0 bg-black bg-opacity-75 flex flex-col justify-center items-center p-4 transition-opacity duration-300">
                                    <h3 className="text-white text-xl font-semibold mb-2">{project.title}</h3>
                                    <p className="text-white text-sm text-center mb-4">{project.description}</p>
                                    <div className="flex space-x-2">
                                        <Button variant="outline" className="text-white border-white hover:bg-white hover:text-black" onClick={() => window.open(project.githubLink, '_blank', 'noopener,noreferrer')}>
                                            <Github className="mr-2 h-4 w-4" /> GitHub
                                        </Button>
                                        <Button variant="outline" className="text-white border-white hover:bg-white hover:text-black" onClick={() => window.open(project.deploymentLink, '_blank', 'noopener,noreferrer')}>
                                            <ExternalLink className="mr-2 h-4 w-4" /> Live Demo
                                        </Button>
                                    </div>
                                </div>
                            )}
                        </div>
                        <CardContent className="p-4">
                            <h3 className="font-semibold text-lg mb-2">{project.title}</h3>
                            <p className="text-sm text-muted-foreground mb-4">
                                {project.description.length > 100
                                    ? `${project.description.substring(0, 100)}...`
                                    : project.description}
                            </p>
                            <div className="flex flex-wrap gap-2">
                                {project.techStack.map((tech, index) => (
                                    <Badge key={index} variant="secondary">{tech}</Badge>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    )
}