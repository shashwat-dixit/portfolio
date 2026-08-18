import React from "react";
import BlurFade from "@/components/magicui/blur-fade";
import BlurFadeText from "@/components/magicui/blur-fade-text";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { DATA } from "@/data/resume";
import AboutSection from "@/components/section/about-section";
import ContactSection from "@/components/section/contact-section";
import EducationSection from "@/components/section/education-section";
import PhotosSection from "@/components/section/photos-section";
import ProjectsSection from "@/components/section/projects-section";
import SkillsSection from "@/components/section/skills-section";
import WorkTimelineSection from "@/components/section/work-timeline-section";
import SshBanner from "@/components/section/ssh-banner";

const BLUR_FADE_DELAY = 0.04;

const sectionComponents: Record<string, React.ReactNode> = {
  about: (
    <section id="about">
      <BlurFade delay={BLUR_FADE_DELAY * 3}>
        <AboutSection />
      </BlurFade>
    </section>
  ),
  work: (
    <section id="work">
      <BlurFade delay={BLUR_FADE_DELAY * 5}>
        <WorkTimelineSection />
      </BlurFade>
    </section>
  ),
  education: (
    <section id="education">
      <BlurFade delay={BLUR_FADE_DELAY * 7}>
        <EducationSection />
      </BlurFade>
    </section>
  ),
  skills: (
    <section id="skills">
      <BlurFade delay={BLUR_FADE_DELAY * 9}>
        <SkillsSection />
      </BlurFade>
    </section>
  ),
  projects: (
    <section id="projects">
      <BlurFade delay={BLUR_FADE_DELAY * 11}>
        <ProjectsSection />
      </BlurFade>
    </section>
  ),
  photos: <PhotosSection />,
  contact: (
    <section id="contact">
      <BlurFade delay={BLUR_FADE_DELAY * 16}>
        <ContactSection />
      </BlurFade>
    </section>
  ),
};

export default function HomePage() {
  const orderedSections = Object.entries(DATA.sections)
    .filter(([, s]) => s.enabled)
    .sort(([, a], [, b]) => a.order - b.order)
    .map(([key]) => key);

  return (
    <main className="min-h-dvh flex flex-col gap-14 relative">
      <section id="hero">
        <div className="mx-auto w-full max-w-2xl space-y-8">
          <div className="gap-2 gap-y-6 flex flex-col md:flex-row justify-between">
            <div className="gap-2 flex flex-col order-2 md:order-1">
              <BlurFadeText
                delay={BLUR_FADE_DELAY}
                className="text-3xl font-semibold tracking-tighter sm:text-4xl lg:text-5xl"
                yOffset={8}
                text={`Hi, I'm ${DATA.name.split(" ")[0]}`}
              />
              <BlurFadeText
                className="text-muted-foreground max-w-[600px] md:text-lg lg:text-xl"
                delay={BLUR_FADE_DELAY}
                text={DATA.description}
              />
            </div>
            <BlurFade delay={BLUR_FADE_DELAY} className="order-1 md:order-2">
              <Avatar className="size-24 md:size-32 border rounded-full shadow-lg ring-4 ring-muted">
                <AvatarImage alt={DATA.name} src={DATA.avatarUrl} />
                <AvatarFallback>{DATA.initials}</AvatarFallback>
              </Avatar>
            </BlurFade>
          </div>
        </div>
      </section>
      <section id="terminal">
        <BlurFade delay={BLUR_FADE_DELAY * 2} className="mx-auto w-full max-w-2xl">
          <SshBanner />
        </BlurFade>
      </section>
      {orderedSections.map((key) => (
        <React.Fragment key={key}>
          {sectionComponents[key]}
        </React.Fragment>
      ))}
    </main>
  );
}
