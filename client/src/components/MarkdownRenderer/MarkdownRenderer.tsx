import { useMemo } from 'react';
import { marked, type MarkedOptions } from 'marked';
import DOMPurify from 'dompurify';

interface Props {
  /** Markdown string to render */
  content: string;

  /** Optional class name for the wrapper element */
  className?: string;
}

function MarkdownRenderer({content, className}: Props) {
    const options: MarkedOptions = {
        async: false,
        breaks: true,
        gfm: true,
        silent: true
    }

    // `useMemo` keeps the expensive parse step from running on every render
    const html = useMemo(() => {
        const raw = marked(content, options);
        return DOMPurify.sanitize(raw as string);
    }, [content]);

    return (
        <div
            className={className}
            dangerouslySetInnerHTML={{ __html: html }}
        />
    );
};

export default MarkdownRenderer
