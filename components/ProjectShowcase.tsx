'use client'
import { useState } from 'react'
import Image from 'next/image'
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
        title: "Jamin AI",
        description: "Jamin is an LLM aggregator.",
        image: "https://github.com/shashwat-dixit/jamin/blob/master/img/jaminFrontPage.png",
        githubLink: "https://github.com/shashwat-dixit/jamin",
        deploymentLink: "https://jamin.shashwatdixit.dev",
        techStack: ["React", "TypeScript", "Langchain"]
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
            <h1 className="text-3xl font-bold mb-12 text-center">My Projects</h1>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-7xl mx-auto">
                {/* Row 1: Large + Small */}
                <Card
                    className="overflow-hidden transition-all duration-300 hover:shadow-xl md:h-[600px]"
                    onMouseEnter={() => setHoveredProject(projects[0].id)}
                    onMouseLeave={() => setHoveredProject(null)}
                >
                    <div className="relative h-full">
                        <Image
                            src={projects[0].image}
                            alt={projects[0].title}
                            className="w-full h-full object-cover"
                        />
                        {hoveredProject === projects[0].id && (
                            <div className="absolute inset-0 bg-black bg-opacity-75 flex flex-col justify-center items-center p-8 transition-opacity duration-300">
                                <h3 className="text-white text-3xl font-semibold mb-6">{projects[0].title}</h3>
                                <p className="text-white text-xl text-center mb-8 max-w-lg">{projects[0].description}</p>
                                <div className="flex flex-wrap gap-3 mb-8 justify-center">
                                    {projects[0].techStack.map((tech, index) => (
                                        <Badge key={index} variant="secondary" className="text-white text-lg px-4 py-1">{tech}</Badge>
                                    ))}
                                </div>
                                <div className="flex space-x-6">
                                    <Button size="lg" variant="outline" className="text-white border-white hover:bg-white hover:text-black"
                                        onClick={() => window.open(projects[0].githubLink, '_blank', 'noopener,noreferrer')}>
                                        <Github className="mr-2 h-6 w-6" /> GitHub
                                    </Button>
                                    <Button size="lg" variant="outline" className="text-white border-white hover:bg-white hover:text-black"
                                        onClick={() => window.open(projects[0].deploymentLink, '_blank', 'noopener,noreferrer')}>
                                        <ExternalLink className="mr-2 h-6 w-6" /> Live Demo
                                    </Button>
                                </div>
                            </div>
                        )}
                    </div>
                </Card>

                <div className="grid grid-rows-2 gap-8">
                    {projects.slice(1, 3).map((project) => (
                        <Card
                            key={project.id}
                            className="overflow-hidden transition-all duration-300 hover:shadow-xl md:h-[280px]"
                            onMouseEnter={() => setHoveredProject(project.id)}
                            onMouseLeave={() => setHoveredProject(null)}
                        >
                            <div className="relative h-full">
                                <Image
                                    src={project.image}
                                    alt={project.title}
                                    className="w-full h-full object-cover"
                                />
                                {hoveredProject === project.id && (
                                    <div className="absolute inset-0 bg-black bg-opacity-75 flex flex-col justify-center items-center p-6 transition-opacity duration-300">
                                        <h3 className="text-white text-2xl font-semibold mb-3">{project.title}</h3>
                                        <p className="text-white text-lg text-center mb-4">{project.description}</p>
                                        <div className="flex flex-wrap gap-2 mb-4 justify-center">
                                            {project.techStack.map((tech, index) => (
                                                <Badge key={index} variant="secondary" className="text-white">{tech}</Badge>
                                            ))}
                                        </div>
                                        <div className="flex space-x-4">
                                            <Button variant="outline" className="text-white border-white hover:bg-white hover:text-black"
                                                onClick={() => window.open(project.githubLink, '_blank', 'noopener,noreferrer')}>
                                                <Github className="mr-2 h-4 w-4" /> GitHub
                                            </Button>
                                            <Button variant="outline" className="text-white border-white hover:bg-white hover:text-black"
                                                onClick={() => window.open(project.deploymentLink, '_blank', 'noopener,noreferrer')}>
                                                <ExternalLink className="mr-2 h-4 w-4" /> Live Demo
                                            </Button>
                                        </div>
                                    </div>
                                )}
                            </div>
                        </Card>
                    ))}
                </div>

                {/* Row 2: Two Equal Sized Cards */}
                {projects.slice(3).map((project) => (
                    <Card
                        key={project.id}
                        className="overflow-hidden transition-all duration-300 hover:shadow-xl md:h-[400px]"
                        onMouseEnter={() => setHoveredProject(project.id)}
                        onMouseLeave={() => setHoveredProject(null)}
                    >
                        <div className="relative h-full">
                            <Image
                                src={project.image}
                                alt={project.title}
                                className="w-full h-full object-cover"
                            />
                            {hoveredProject === project.id && (
                                <div className="absolute inset-0 bg-black bg-opacity-75 flex flex-col justify-center items-center p-6 transition-opacity duration-300">
                                    <h3 className="text-white text-2xl font-semibold mb-4">{project.title}</h3>
                                    <p className="text-white text-lg text-center mb-6">{project.description}</p>
                                    <div className="flex flex-wrap gap-2 mb-6 justify-center">
                                        {project.techStack.map((tech, index) => (
                                            <Badge key={index} variant="secondary" className="text-white">{tech}</Badge>
                                        ))}
                                    </div>
                                    <div className="flex space-x-4">
                                        <Button variant="outline" className="text-white border-white hover:bg-white hover:text-black"
                                            onClick={() => window.open(project.githubLink, '_blank', 'noopener,noreferrer')}>
                                            <Github className="mr-2 h-5 w-5" /> GitHub
                                        </Button>
                                        <Button variant="outline" className="text-white border-white hover:bg-white hover:text-black"
                                            onClick={() => window.open(project.deploymentLink, '_blank', 'noopener,noreferrer')}>
                                            <ExternalLink className="mr-2 h-5 w-5" /> Live Demo
                                        </Button>
                                    </div>
                                </div>
                            )}
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    )
}