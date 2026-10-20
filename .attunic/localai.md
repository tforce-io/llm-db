# LocalAI Guideline

This guide state the facts about the yaml files in `yaml/LocalAI` folders.

- Configurations are optimized for 48G VRAM system.
- Parameters should follow recommendation from model developer, not convention from other files.

## Text Generation Models

- Create 2 variants for model that support vision: base and mmproj.
- Use Qwen3.6-27B as reference for file formatting.
- Use Q4_K_M, Q8_0 variants for Dense models.
- Prefer Q4_K_XL, Q8_K_XL variants for MoE models.
- Prefer F16 variant over BF16 for MMProj models. Do not use Q8 or F32.
- Always verify all information against model card: max context length, capabilities, modalities, structured output, templating support,...
- Use sampling parameters for agentic coding task, based on recommendation of model developer.

## Image Generation Models

- Create 2 variants for model: base and unrestricted.
- Use Qwen-Image-2.1-7B as reference for file formatting.
- Always verify all information against model card: capabilities, modalities,...
