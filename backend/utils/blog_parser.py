# Parses Blog Metadata & Markdown Content

# TODO: Consider parsing the blog content further to optimize images. 

from typing import Any, Dict, Tuple
import yaml

def parse_blog(file_path: str) -> Tuple[Dict[str, Any], str]:
    """
    Parse YAML front matter bounded by the first pair of lines that are 
    exactly '---' at the very start of the file. All subsequent '---' stay in 
    the body.

    Returns (metadata, body).
    """
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    # Split into lines but keep line endings so we preserve the body unchanged
    lines = content.splitlines(keepends=True)
    if not lines or lines[0].strip() != '---':
        # No front matter at the very top
        return {}, content

    # Find the terminating line that's exactly '---'
    end_idx = None
    for i in range(1, len(lines)):
        token = lines[i].strip()
        if token == '---':
            end_idx = i
            break

    if end_idx is None:
        # No closing marker — treat whole file as body
        return {}, content

    # Extract YAML block (between the markers) and the remaining body
    yaml_block = ''.join(lines[1:end_idx])
    # Preserve body exactly, later '---' untouched
    body = ''.join(lines[end_idx + 1:])

    # Parse YAML safely; on error, fall back to empty metadata
    try:
        data = yaml.safe_load(yaml_block.strip())
        metadata: Dict[str, Any] = data if isinstance(data, dict) else {}

    except yaml.YAMLError:
        metadata = {}

    return metadata, body