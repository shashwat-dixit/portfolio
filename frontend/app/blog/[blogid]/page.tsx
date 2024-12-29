type Props = {
  params: {
    blogid: string
  }
}

export default function BlogPage({ params }: Props) {
  return (
    <div className="flex min-h-screen flex-col items-center justify-between p-24">
      <h1>Blog Post: {params.blogid}</h1>
    </div>
  );
}