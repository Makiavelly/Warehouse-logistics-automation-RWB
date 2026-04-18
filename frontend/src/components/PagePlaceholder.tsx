interface PagePlaceholderProps {
  minHeight?: string | number;
}

export function PagePlaceholder(props: PagePlaceholderProps) {
  return <div style={{ minHeight: typeof props.minHeight === "undefined" ? "60vh" : props.minHeight }} />;
}
