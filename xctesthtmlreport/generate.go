package xctesthtmlreport

func (b *BitriseXchtmlGenerator) Generate(outputPath, xcresultPath string) error {
	b.logger.Printf("Generating report for: %s", xcresultPath)
	params := []string{"--output", outputPath, xcresultPath}
	cmd := b.commandFactory.Create(b.toolPath, params, nil)
	_, err := cmd.RunAndReturnTrimmedCombinedOutput()
	return err
}
