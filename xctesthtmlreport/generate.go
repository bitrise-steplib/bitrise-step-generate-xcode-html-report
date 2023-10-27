package xctesthtmlreport

func (b *BitriseXchtmlGenerator) Generate(outputPath, xcresultPath string) error {
	b.Logger.Printf("Generating report for: %s", xcresultPath)
	params := []string{"--output", outputPath, xcresultPath}
	cmd := b.CommandFactory.Create(b.toolPath, params, nil)
	_, err := cmd.RunAndReturnTrimmedCombinedOutput()
	return err
}
